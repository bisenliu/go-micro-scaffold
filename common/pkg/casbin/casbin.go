package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	entadapter "github.com/casbin/ent-adapter"
	"go.uber.org/zap"

	"common/databases/rdbms"
)

// EnforcerParams 定义了创建Casbin Enforcer所需的依赖
type EnforcerParams struct {
	Client *rdbms.Client
	Logger *zap.Logger
}

// NewEnforcer 创建一个 Casbin SyncedCachedEnforcer 实例
func NewEnforcer(client *rdbms.Client, logger *zap.Logger) (*casbin.SyncedCachedEnforcer, error) {
	// 使用数据源名称创建适配器，复用现有的数据库配置
	// 获取数据库配置信息
	config := client.Config()

	// 构建数据源名称
	charset := config.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		charset)

	logger.Info("Initializing Casbin adapter",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.String("database", config.Database))

	// 创建适配器
	// 注意：Casbin表结构现在通过Ent的schema管理，使用ent迁移功能创建
	a, err := entadapter.NewAdapter("mysql", dataSourceName)
	if err != nil {
		logger.Error("Failed to create casbin ent adapter", zap.Error(err))
		return nil, fmt.Errorf("failed to create casbin ent adapter: %w", err)
	}

	// 从字符串加载 Casbin 模型
	text := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
	`
	m, err := model.NewModelFromString(text)
	if err != nil {
		logger.Error("Failed to load casbin model from string", zap.Error(err))
		return nil, fmt.Errorf("failed to load casbin model from string: %w", err)
	}

	// 创建 Enforcer 实例
	enforcer, err := casbin.NewSyncedCachedEnforcer(m, a)
	if err != nil {
		logger.Error("Failed to create casbin enforcer", zap.Error(err))
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// 设置缓存过期时间（例如：10分钟）
	enforcer.SetExpireTime(10 * 60)

	// 从数据库加载策略
	logger.Info("Loading casbin policy from database")
	if err := enforcer.LoadPolicy(); err != nil {
		logger.Error("Failed to load casbin policy", zap.Error(err))
		return nil, fmt.Errorf("failed to load casbin policy: %w", err)
	}

	logger.Info("Casbin enforcer initialized successfully")
	return enforcer, nil
}
