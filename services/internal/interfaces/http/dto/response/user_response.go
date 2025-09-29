package response

import (
	commonResponse "common/response"
	"services/internal/domain/user/entity"
)

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID          string `json:"id"`
	OpenID      string `json:"open_id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []*UserInfoResponse `json:"users"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// ToUserInfoResponse 将用户实体转换为用户信息响应
func ToUserInfoResponse(user *entity.User) *UserInfoResponse {
	if user == nil {
		return nil
	}

	response := &UserInfoResponse{
		ID:          user.ID(),
		OpenID:      user.OpenID(),
		Name:        user.Name(),
		Gender:      user.Gender(),
		PhoneNumber: user.PhoneNumber(),
		CreatedAt:   user.GetCreatedAt(),
		UpdatedAt:   user.GetUpdatedAt(),
	}

	// 可以添加更多字段的转换逻辑

	return response
}

// ToUserListResponse 将用户实体列表转换为用户列表响应
func ToUserListResponse(users []*entity.User, total int64, page, pageSize int) ([]*UserInfoResponse, *commonResponse.Pagination) {
	userResponses := make([]*UserInfoResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, ToUserInfoResponse(user))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	pagination := &commonResponse.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return userResponses, pagination
}
