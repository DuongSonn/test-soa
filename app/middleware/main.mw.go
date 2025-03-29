package middleware

import (
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/redis"
)

type MiddlewareCollections struct {
	RateLimitMw ICustomMiddleware
	AuthMw      ICustomMiddleware
	AdminMw     ICustomMiddleware
}

func RegisterMiddleware(
	redisClient redis.IRedisClient,
	postgresRepo repository.RepositoryCollections,
	helpers helper.HelperCollections,
) MiddlewareCollections {
	return MiddlewareCollections{
		RateLimitMw: NewRateLimitMiddleware(redisClient),
		AuthMw:      NewAuthMiddleware(postgresRepo, helpers),
		AdminMw:     NewAdminMiddleware(),
	}
}
