package model

import (
	"sondth-test_soa/app/entity"

	"github.com/google/uuid"
)

// UserLoginRequest struct
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserRegisterRequest struct
type UserRegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	Fullname        string `json:"fullname" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
type UserRegisterResponse struct{}

// UserForgotPasswordRequest struct
type ChangeUserPasswordRequest struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	OldPassword     string    `json:"old_password" validate:"required"`
	NewPassword     string    `json:"new_password" validate:"required"`
	ConfirmPassword string    `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
type ChangeUserPasswordResponse struct{}

// UpdateUserRequest struct
type UpdateUserRequest struct {
	ID       uuid.UUID `param:"id" validate:"required"`
	Username *string   `json:"username"`
	Fullname *string   `json:"fullname"`
	Role     *string   `json:"role"`
}
type UpdateUserResponse struct {
	User entity.User `json:"user"`
}

// ForgetPasswordRequest struct
type ForgetPasswordRequest struct {
	Username        string `json:"username" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
type ForgetPasswordResponse struct{}

// GetUserRequest struct
type GetUsersRequest struct {
	Name  *string `json:"name"`
	Role  *string `json:"role"`
	Page  *int    `json:"page"`
	Limit *int    `json:"limit"`
}
type GetUsersResponse struct {
	Users []entity.User `json:"users"`
	Count int64         `json:"count"`
}
