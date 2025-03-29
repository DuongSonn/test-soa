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
	"sondth-test_soa/utils"
)

type userHandler struct {
	services service.ServiceCollections
	mws      middleware.MiddlewareCollections
}

func NewUserControllerV1(router *gin.Engine, services service.ServiceCollections, mws middleware.MiddlewareCollections) {
	handler := userHandler{services, mws}

	group := router.Group("api/v1/user")
	{
		group.POST("/register", handler.register)
		group.POST("/login", handler.login)
		group.POST("/forget-password", handler.forgetPassword)
		group.POST("/change-password", handler.changePassword)
		group.POST("/update/:id", handler.updateUser)
		group.POST("/list", handler.getUsers, mws.AdminMw.Handler())
	}
}

func (h *userHandler) register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.UserService.Register(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *userHandler) login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.UserService.Login(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *userHandler) forgetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.ForgetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.UserService.ForgetPassword(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *userHandler) changePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.ChangeUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.UserService.ChangePassword(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *userHandler) updateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.UpdateUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.UserService.UpdateUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *userHandler) getUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req model.GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
		return
	}

	res, err := h.services.UserService.GetUsers(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
