package validation

import (
	"common/databases/mysql"
	"services/internal/domain/user/repository"

	"github.com/go-playground/validator/v10"
)

// UserInfrastructureValidator 用户基础设施验证器
// 负责数据库层面的验证，如唯一性检查等
type UserInfrastructureValidator struct {
	db       *mysql.EntClient
	userRepo repository.UserRepository
	validate *validator.Validate
}

// NewUserInfrastructureValidator 创建用户基础设施验证器
func NewUserInfrastructureValidator(db *mysql.EntClient, userRepo repository.UserRepository) *UserInfrastructureValidator {
	// 创建验证器实例
	validate := validator.New()

	// // 注册自定义验证规则
	// validation.RegisterCustomValidations(validate)

	// // 注册标签名函数
	// validation.RegisterTagNameFunc(validate)

	return &UserInfrastructureValidator{
		db:       db,
		userRepo: userRepo,
		validate: validate,
	}
}
