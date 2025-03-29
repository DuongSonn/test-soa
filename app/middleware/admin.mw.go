package middleware

import (
	"net/http"
	"sondth-test_soa/app/entity"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"

	"github.com/gin-gonic/gin"
)

type adminMiddleware struct {
}

func NewAdminMiddleware() ICustomMiddleware {
	return &adminMiddleware{}
}

func (m *adminMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get(string(utils.USER_CONTEXT_KEY))
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeUnauthorized)))
			c.Abort()
			return
		}
		userEntity, ok := user.(*entity.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeUnauthorized)))
			c.Abort()
			return
		}
		if !userEntity.IsAdmin() {
			c.JSON(http.StatusForbidden, errors.FormatErrorResponse(errors.New(errors.ErrCodeForbidden)))
			c.Abort()
			return
		}

		c.Next()
	}
}
