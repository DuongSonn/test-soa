package errors

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"

	logger "sondth-test_soa/package/log"
	_validator "sondth-test_soa/package/validator"
	"sondth-test_soa/utils"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	LangVN = "vi"
	LangEN = "en"
)

const (
	ErrCodeValidatorRequired     = 1
	ErrCodeValidatorFormat       = 2
	ErrCodeValidatorVerifiedData = 3

	// User Error
	ErrCodeUserNotFound     = 10
	ErrCodeUserExisted      = 11
	ErrCodeWishlistNotFound = 12

	// OAuth Error
	ErrCodeTokenExpired      = 20
	ErrCodeIncorrectPassword = 21

	// Category Error
	ErrCodeCategoryExisted  = 30
	ErrCodeCategoryNotFound = 31

	// Product Error
	ErrCodeProductNotFound          = 40
	ErrCodeProductExisted           = 41
	ErrCodeProductAlreadyInWishlist = 42
	ErrCodeReviewNotFound          = 43
	ErrCodeReviewAlreadyExists     = 44

	// System Error
	ErrCodeInternalServerError = 500
	ErrCodeTimeout             = 408
	ErrCodeForbidden           = 403
	ErrCodeUnauthorized        = 402
	ErrCodeRateLimitExceeded   = 429
)

var messages = map[int]map[string]string{
	// Validator
	ErrCodeValidatorRequired: {
		LangVN: "%s không được bỏ trống. Vui lòng kiểm tra lại",
		LangEN: "%s is required. Please check again",
	},
	ErrCodeValidatorFormat: {
		LangVN: "%s không hợp lệ. Vui lòng kiểm tra lại",
		LangEN: "%s is invalid. Please check again",
	},
	ErrCodeValidatorVerifiedData: {
		LangVN: "%s không chính xác. Vui lòng kiểm tra lại",
		LangEN: "%s is incorrect. Please check again",
	},

	// System Error
	ErrCodeInternalServerError: {
		LangVN: "Lỗi hệ thống",
		LangEN: "Internal server error",
	},
	ErrCodeTimeout: {
		LangVN: "Hết thời gian xử lý",
		LangEN: "Request timeout",
	},
	ErrCodeForbidden: {
		LangVN: "Không có quyền truy cập",
		LangEN: "Access forbidden",
	},
	ErrCodeUnauthorized: {
		LangVN: "Chưa xác thực",
		LangEN: "Unauthorized",
	},
	ErrCodeRateLimitExceeded: {
		LangVN: "Đã vượt quá giới hạn yêu cầu",
		LangEN: "Rate limit exceeded",
	},

	// User Error
	ErrCodeUserNotFound: {
		LangVN: "Không tìm thấy người dùng. Vui lòng kiểm tra lại",
		LangEN: "User not found. Please check again",
	},
	ErrCodeUserExisted: {
		LangVN: "Người dùng đã đăng ký tài khoản. Vui lòng kiểm tra lại",
		LangEN: "User already exists. Please check again",
	},
	ErrCodeWishlistNotFound: {
		LangVN: "Không tìm thấy danh sách yêu thích. Vui lòng kiểm tra lại",
		LangEN: "Wishlist not found. Please check again",
	},

	// OAuth Error
	ErrCodeTokenExpired: {
		LangVN: "Token đã hết hạn",
		LangEN: "Token has expired",
	},
	ErrCodeIncorrectPassword: {
		LangVN: "Mật khẩu không chính xác",
		LangEN: "Password is incorrect",
	},

	// Category Error
	ErrCodeCategoryExisted: {
		LangVN: "Danh mục đã tồn tại. Vui lòng kiểm tra lại",
		LangEN: "Category already exists. Please check again",
	},
	ErrCodeCategoryNotFound: {
		LangVN: "Danh mục không tồn tại. Vui lòng kiểm tra lại",
		LangEN: "Category not found. Please check again",
	},

	// Product Error
	ErrCodeProductNotFound: {
		LangVN: "Sản phẩm không tồn tại. Vui lòng kiểm tra lại",
		LangEN: "Product not found. Please check again",
	},
	ErrCodeProductExisted: {
		LangVN: "Sản phẩm đã tồn tại. Vui lòng kiểm tra lại",
		LangEN: "Product already exists. Please check again",
	},
	ErrCodeProductAlreadyInWishlist: {
		LangVN: "Sản phẩm đã tồn tại trong danh sách yêu thích. Vui lòng kiểm tra lại",
		LangEN: "Product already exists in wishlist. Please check again",
	},
	ErrCodeReviewNotFound: {
		LangVN: "Không tìm thấy đánh giá. Vui lòng kiểm tra lại",
		LangEN: "Review not found. Please check again",
	},
	ErrCodeReviewAlreadyExists: {
		LangVN: "Bạn đã đánh giá sản phẩm này. Vui lòng kiểm tra lại",
		LangEN: "You have already reviewed this product. Please check again",
	},
}

func New(code int) *CustomError {
	return &CustomError{
		Code:    code,
		Message: GetMessage(code),
	}
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func NewValidatorError(err error) *CustomError {
	logger.WithCtx(context.Background()).Info("Validator Error", slog.String("error", err.Error()))
	var validatorErr validator.ValidationErrors
	if errors.As(err, &validatorErr) {
		errDetail := validatorErr[0]

		field := errDetail.Field()
		tag := errDetail.Tag()

		code := convertValidatorTag(tag)
		return &CustomError{
			Code:    code,
			Message: GetCustomMessage(code, field),
		}

	}

	return New(ErrCodeInternalServerError)
}

func GetCustomMessage(code int, args ...any) string {
	msg, ok := messages[code][LangEN]
	if !ok {
		return messages[ErrCodeInternalServerError][LangEN]
	}

	return fmt.Sprintf(msg, args...)
}

func GetMessage(code int) string {
	msg, ok := messages[code][LangEN]
	if !ok {
		return messages[ErrCodeInternalServerError][LangEN]
	}

	return msg
}

func FormatErrorResponse(err error) utils.HttpResponse {
	if e, ok := err.(*CustomError); ok {
		return utils.HttpResponse{
			Success: false,
			Data:    e,
		}
	}

	return utils.HttpResponse{
		Success: false,
		Data:    New(ErrCodeInternalServerError),
	}
}

func (err *CustomError) Error() string {
	return err.Message
}

func (err *CustomError) GetCode() int {
	return err.Code
}

// --------------------------------------
func convertValidatorTag(tag string) int {
	switch tag {
	case _validator.EMAIL, _validator.PHONE_NUMBER:
		return ErrCodeValidatorFormat
	case _validator.EQUAL_FIELD:
		return ErrCodeValidatorVerifiedData
	default:
		return ErrCodeValidatorRequired
	}
}
