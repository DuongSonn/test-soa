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

type categoryHandler struct {
	services service.ServiceCollections
	mws      middleware.MiddlewareCollections
}

func NewCategoryControllerV1(router *gin.Engine, services service.ServiceCollections, mws middleware.MiddlewareCollections) {
	handler := categoryHandler{services, mws}

	group := router.Group("api/v1/category")
	{
		adminGroup := group.Group("/", mws.AdminMw.Handler())
		{
			adminGroup.POST("/create", handler.create)
			adminGroup.GET("/summary", handler.getCategoriesSummary)
		}

		group.POST("/list", handler.getCategories)
	}
}

func (h *categoryHandler) create(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	res, err := h.services.CategorySvc.Create(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, utils.FormatSuccessResponse(res))
}

func (h *categoryHandler) getCategories(c *gin.Context) {
	var req model.GetCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	resp, err := h.services.CategorySvc.GetCategories(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(resp))
}

func (h *categoryHandler) getCategoriesSummary(c *gin.Context) {
	var req model.GetCategoriesSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resErr := errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, errors.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	resp, err := h.services.CategorySvc.GetCategoriesSummary(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(resp))
}
