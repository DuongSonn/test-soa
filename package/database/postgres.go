package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sondth-test_soa/config"
	"sondth-test_soa/utils"
)

type PostgresClient struct {
	db *gorm.DB
}

func NewPostgresClient(conf config.Configuration) (*PostgresClient, error) {
	// DB logging config
	logLevel := logger.Info
	if conf.Server.Mode == utils.RELEASE_MODE {
		logLevel = logger.Silent
	}
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logLevel,               // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Enable color
		},
	)

	// DB connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.PostgresDB.Host,
		conf.PostgresDB.User,
		conf.PostgresDB.Password,
		conf.PostgresDB.Database,
		conf.PostgresDB.Port,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting database: %v", err)
	}
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	client := &PostgresClient{db: conn}

	return client, nil
}

// GetDB returns the underlying *gorm.DB instance
func (p *PostgresClient) GetDB() *gorm.DB {
	return p.db
}

// BeginTransaction starts a new transaction
func (p *PostgresClient) BeginTransaction() *gorm.DB {
	return p.db.Begin()
}
