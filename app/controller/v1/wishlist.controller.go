package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"sondth-test_soa/app/model"
	"sondth-test_soa/app/service"
	"sondth-test_soa/package/errors"
)

type wishlistHandler struct {
	services service.ServiceCollections
}

func NewWishlistControllerV1(router *gin.Engine, services service.ServiceCollections) {
	handler := wishlistHandler{services}

	group := router.Group("api/v1/wishlist")
	{
		group.POST("/add", handler.addToWishlist)
		group.POST("/delete/:id", handler.deleteFromWishlist)
		group.POST("/list", handler.getWishlists)
		group.GET("/summary", handler.getWishlistsSummary)
	}
}

func (h *wishlistHandler) addToWishlist(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.WishlistSvc.AddToWishlist(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *wishlistHandler) deleteFromWishlist(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.DeleteFromWishlistRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.WishlistSvc.DeleteFromWishlist(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *wishlistHandler) getWishlists(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.GetWishlistsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.WishlistSvc.GetWishlists(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *wishlistHandler) getWishlistsSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.GetWishlistsSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.WishlistSvc.GetWishlistsSummary(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
