package middleware

import (
	_errors "errors"
	"net/http"
	"slices"
	"strings"

	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"
	logger "sondth-test_soa/package/log"
	"sondth-test_soa/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	WHITE_LIST_API = []string{
		"/api/v1/user/login",
		"/api/v1/user/register",
	}
)

type authMiddleware struct {
	postgresRepo repository.RepositoryCollections
	helpers      helper.HelperCollections
}

func NewAuthMiddleware(postgresRepo repository.RepositoryCollections, helpers helper.HelperCollections) ICustomMiddleware {
	return &authMiddleware{postgresRepo: postgresRepo, helpers: helpers}
}

func (m *authMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(WHITE_LIST_API, c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeUnauthorized)))
			c.Abort()
			return
		}
		authArr := strings.Split(authHeader, " ")
		if len(authArr) != 2 || authArr[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeUnauthorized)))
			c.Abort()
			return
		}

		payload, err := m.helpers.OAuthHelper.VerifyAccessToken(authArr[1])
		if err != nil {
			if _errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeTokenExpired)))
				c.Abort()
				return
			}

			logger.WithCtx(c).Error("VerifyAccessToken", err)
			c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeUnauthorized)))
			c.Abort()
			return
		}

		user, err := m.postgresRepo.UserRepo.FindOneByFilter(c, nil, &repository.FindUserByFilter{
			ID: &payload.UserID,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, errors.FormatErrorResponse(errors.New(errors.ErrCodeUnauthorized)))
			c.Abort()
			return
		}

		c.Set(string(utils.USER_CONTEXT_KEY), user)
		c.Next()
	}
}
