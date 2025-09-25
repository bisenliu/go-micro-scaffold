package response

import (
	"time"
)

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ValidationErrorResponse 验证错误响应
type ValidationErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
}

// UserResponse 用户响应DTO
type UserResponse struct{}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID        int64     `json:"id"`
	OpenID    string    `json:"open_id"`
	UnionID   string    `json:"union_id"`
	Phone     string    `json:"phone"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Gender    string    `json:"gender"`
	Birthday  *string   `json:"birthday"`
	Location  string    `json:"location"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []*UserInfoResponse `json:"users"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
	}
}

// ErrorResponseWithValidation 带验证错误的错误响应
func ErrorResponseWithValidation(code int, message string, errors map[string]string) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}
