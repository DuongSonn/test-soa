package helper

import (
	"encoding/json"
	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/model"
	"sondth-test_soa/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type oAuthHelper struct {
	config config.Configuration
}

func NewOAuthHelper(
	config config.Configuration,
) IOAuthHelper {
	return &oAuthHelper{config: config}
}

func (h *oAuthHelper) GenerateAccessToken(user entity.User) (string, error) {
	payload := &model.UserJWTPayload{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    h.config.Jwt.Issuer,
		},
	}
	accessToken, err := h.GenerateToken(payload, h.config.Jwt.UserAccessTokenKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (h *oAuthHelper) GenerateRefreshToken(user entity.User) (string, error) {
	payload := &model.UserJWTPayload{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 30)),
			Issuer:    h.config.Jwt.Issuer,
		},
	}
	refreshToken, err := h.GenerateToken(payload, h.config.Jwt.UserRefreshTokenKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (h *oAuthHelper) VerifyAccessToken(tokenString string) (*model.UserJWTPayload, error) {
	token, err := h.VerifyToken(tokenString, h.config.Jwt.UserAccessTokenKey)
	if err != nil {
		return nil, err
	}

	// Convert map claims to JSON bytes
	jsonBytes, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, err
	}

	// Unmarshal into UserJWTPayload
	var payload model.UserJWTPayload
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (h *oAuthHelper) VerifyRefreshToken(tokenString string) (*model.UserJWTPayload, error) {
	token, err := h.VerifyToken(tokenString, h.config.Jwt.UserRefreshTokenKey)
	if err != nil {
		return nil, err
	}

	// Convert map claims to JSON bytes
	jsonBytes, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, err
	}

	// Unmarshal into UserJWTPayload
	var payload model.UserJWTPayload
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (h *oAuthHelper) GenerateToken(claims jwt.Claims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func (h *oAuthHelper) VerifyToken(tokenString string, key string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return token, nil
}
