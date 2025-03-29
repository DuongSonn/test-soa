package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sondth-test_soa/package/errors"
	"sondth-test_soa/package/redis"
)

const (
	MAX_REQUESTS = 100
)

type rateLimitMiddleware struct {
	redisClient redis.IRedisClient
}

func NewRateLimitMiddleware(redisClient redis.IRedisClient) *rateLimitMiddleware {
	return &rateLimitMiddleware{
		redisClient: redisClient,
	}
}

func (mw *rateLimitMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = c.Request.RemoteAddr
		}
		if clientIP == "" {
			clientIP = c.GetHeader("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = c.GetHeader("X-Forwarded-For")
		}
		if clientIP == "" {
			c.JSON(http.StatusTooManyRequests, errors.FormatErrorResponse(errors.New(errors.ErrCodeRateLimitExceeded)))
			c.Abort()
			return
		}

		key := fmt.Sprintf("ip:%s", clientIP)
		value, err := mw.redisClient.Incr(context.Background(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.FormatErrorResponse(errors.New(errors.ErrCodeInternalServerError)))
			c.Abort()
			return
		}
		if value > MAX_REQUESTS {
			c.JSON(http.StatusTooManyRequests, errors.FormatErrorResponse(errors.New(errors.ErrCodeRateLimitExceeded)))
			c.Abort()
			return
		}
		c.Next()
	}
}
