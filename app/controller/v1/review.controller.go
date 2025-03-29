package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"sondth-test_soa/app/middleware"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/service"
	"sondth-test_soa/package/errors"
)

type reviewHandler struct {
	services service.ServiceCollections
	mws      middleware.MiddlewareCollections
}

func NewReviewControllerV1(router *gin.Engine, services service.ServiceCollections, mws middleware.MiddlewareCollections) {
	handler := reviewHandler{services, mws}

	group := router.Group("api/v1/review")
	{
		group.POST("/create", handler.create)
		group.POST("/delete", handler.delete)
		group.POST("/list", handler.getReviews)
		group.GET("/summary", handler.getReviewsSummary)
	}
}

func (h *reviewHandler) create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.ReviewSvc.Create(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *reviewHandler) delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.DeleteReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.ReviewSvc.Delete(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *reviewHandler) getReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.GetReviewsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.ReviewSvc.GetReviews(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *reviewHandler) getReviewsSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.GetReviewsSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.ReviewSvc.GetReviewsSummary(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
