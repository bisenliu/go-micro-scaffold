package response

import (
	"user-services/internal/domain/user/entity"
)

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID          string `json:"id" example:"user_123456789"`        // 用户唯一标识ID
	OpenID      string `json:"open_id" example:"wx_123456789"`     // 第三方平台的唯一标识
	Name        string `json:"name" example:"张三"`                  // 用户姓名
	Gender      int    `json:"gender" example:"200"`               // 性别：100-男性，200-女性，300-其他
	PhoneNumber string `json:"phone_number" example:"13800138000"` // 手机号码
	CreatedAt   int64  `json:"created_at" example:"1640995200000"` // 创建时间戳（毫秒）
	UpdatedAt   int64  `json:"updated_at" example:"1640995200000"` // 更新时间戳（毫秒）
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []*UserInfoResponse `json:"users"`                    // 用户列表
	Total      int64               `json:"total" example:"100"`      // 总记录数
	Page       int                 `json:"page" example:"1"`         // 当前页码
	PageSize   int                 `json:"page_size" example:"10"`   // 每页记录数
	TotalPages int                 `json:"total_pages" example:"10"` // 总页数
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
func ToUserListResponse(users []*entity.User) []*UserInfoResponse {
	userResponses := make([]*UserInfoResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, ToUserInfoResponse(user))
	}

	return userResponses
}
