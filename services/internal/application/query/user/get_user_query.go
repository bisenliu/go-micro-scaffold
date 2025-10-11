package user

import "github.com/go-playground/validator/v10"

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	ID string `json:"id" validate:"required,uuid4"` // 用户ID
}

// Validate 验证查询参数
func (q *GetUserQuery) Validate(validate *validator.Validate) error {
	return validate.Struct(q)
}
