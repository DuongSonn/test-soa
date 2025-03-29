package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"sondth-test_soa/app/controller"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/middleware"
	"sondth-test_soa/app/repository/postgres"
	"sondth-test_soa/app/service"
	"sondth-test_soa/config"
	"sondth-test_soa/package/database"
	"sondth-test_soa/package/redis"
	_validator "sondth-test_soa/package/validator"
	"sondth-test_soa/utils"
)

func main() {
	// Initialize configuration
	rootDir := utils.GetFileRootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)
	configClient, err := config.NewConfigClient(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	conf := configClient.Get()

	// Register repositories
	postgresDB, err := database.NewPostgresClient(conf)
	if err != nil {
		log.Fatalf("Failed to initialize Postgres client: %v", err)
	}
	postgresRepo := postgres.RegisterPostgresRepositories(postgresDB.GetDB())

	// Register redis
	redisClient, err := redis.NewRedisClient(conf)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// Register Others
	helpers := helper.RegisterHelpers(postgresRepo, conf)
	services := service.RegisterServices(helpers, postgresRepo)
	mws := middleware.RegisterMiddleware(redisClient, postgresRepo, helpers)

	// Start HTTP Server
	srv := initHTTPServer(conf, services, mws)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 3 seconds.")
	}
	log.Println("Server exiting")
}

func initHTTPServer(conf config.Configuration, services service.ServiceCollections, mws middleware.MiddlewareCollections) *http.Server {
	// Run gin server
	gin.SetMode(conf.Server.Mode)
	app := gin.Default()

	// Register validator
	v := binding.Validator.Engine().(*validator.Validate)
	v.SetTagName("validate")
	_validator.RegisterCustomValidators(v)

	// Register middleware
	app.Use(mws.RateLimitMw.Handler())
	app.Use(mws.AuthMw.Handler())

	// Register controllers
	app.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	controller.RegisterControllers(app, services, mws)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: app,
	}

	return srv
}
