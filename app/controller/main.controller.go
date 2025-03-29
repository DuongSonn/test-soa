package controller

import (
	"github.com/gin-gonic/gin"

	v1 "sondth-test_soa/app/controller/v1"
	"sondth-test_soa/app/middleware"
	"sondth-test_soa/app/service"
)

func RegisterControllers(router *gin.Engine, services service.ServiceCollections, mws middleware.MiddlewareCollections) {
	v1.NewCategoryControllerV1(router, services, mws)
	v1.NewReviewControllerV1(router, services, mws)
	v1.NewProductControllerV1(router, services, mws)
	v1.NewUserControllerV1(router, services, mws)
	v1.NewWishlistControllerV1(router, services)
}
