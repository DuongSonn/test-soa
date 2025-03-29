package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"sondth-test_soa/app/middleware"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/service"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"
)

type productHandler struct {
	services service.ServiceCollections
	mws      middleware.MiddlewareCollections
}

func NewProductControllerV1(router *gin.Engine, services service.ServiceCollections, mws middleware.MiddlewareCollections) {
	handler := productHandler{services, mws}

	group := router.Group("api/v1/product")
	{
		adminGroup := group.Group("/", mws.AdminMw.Handler())
		{
			adminGroup.POST("/create", handler.create)
			adminGroup.POST("/update", handler.update)
			adminGroup.POST("/delete", handler.delete)
		}

		group.POST("/list", handler.getProducts)
	}
}

func (h *productHandler) create(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	res, err := h.services.ProductSvc.Create(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, utils.FormatSuccessResponse(res))
}

func (h *productHandler) update(c *gin.Context) {
	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	resp, err := h.services.ProductSvc.Update(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(resp))
}

func (h *productHandler) delete(c *gin.Context) {
	var req model.DeleteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	resp, err := h.services.ProductSvc.Delete(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(resp))
}

func (h *productHandler) getProducts(c *gin.Context) {
	var req model.GetProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	resp, err := h.services.ProductSvc.GetProducts(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(resp))
}
