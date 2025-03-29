package middleware

import (
	"github.com/gin-gonic/gin"
)

type ICustomMiddleware interface {
	Handler() gin.HandlerFunc
}
