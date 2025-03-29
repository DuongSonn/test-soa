package helper

import (
	"context"
	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ICategoryHelper interface {
	ValidateCategoryID(ctx context.Context, categoryID uuid.UUID) (*entity.Category, error)
}

type IProductHelper interface {
	ValidateProductID(ctx context.Context, productID uuid.UUID) (*entity.Product, error)
}

type IOAuthHelper interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)
	VerifyAccessToken(tokenString string) (*model.UserJWTPayload, error)
	VerifyRefreshToken(tokenString string) (*model.UserJWTPayload, error)

	VerifyToken(tokenString string, key string) (*jwt.Token, error)
	GenerateToken(claims jwt.Claims, key string) (string, error)
}

type IUserHelper interface {
	IsValidRole(role string) bool
	IsValidPermission(ctx context.Context, userID uuid.UUID) bool
}
