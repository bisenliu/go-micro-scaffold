package mysql

import (
	"database/sql"
	"fmt"
	"sync"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
)

// Manager 数据库管理器
type Manager struct {
	clients     sync.Map
	config      *config.Config
	logger      *zap.Logger
	aliases     map[string]string // 别名到数据库名的映射
	defaultName string            // 默认数据库名称
}

// ManagerParams 管理器依赖参数
type ManagerParams struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

// NewManager 创建数据库管理器
func NewManager(params ManagerParams) (*Manager, error) {
	manager := &Manager{
		config:  params.Config,
		logger:  params.Logger,
		aliases: make(map[string]string),
	}

	// 加载别名配置
	if err := manager.loadAliases(); err != nil {
		return nil, fmt.Errorf("failed to load database aliases: %w", err)
	}

	// 初始化所有客户端
	if err := manager.initializeClients(); err != nil {
		return nil, err
	}

	return manager, nil
}

// loadAliases 加载数据库别名配置
func (m *Manager) loadAliases() error {
	if m.config.DatabaseAliases != nil {
		for alias, dbName := range m.config.DatabaseAliases {
			// 验证别名指向的数据库是否存在
			if _, exists := m.config.Databases[dbName]; !exists {
				return fmt.Errorf("database alias '%s' points to non-existent database '%s'", alias, dbName)
			}
			m.aliases[alias] = dbName
		}
	}

	// 设置默认数据库
	if defaultDB, exists := m.aliases["default"]; exists {
		m.defaultName = defaultDB
	} else if len(m.config.Databases) > 0 {
		// 如果没有配置默认别名，使用第一个数据库作为默认
		for name := range m.config.Databases {
			m.defaultName = name
			break
		}
	}

	return nil
}

// initializeClients 初始化所有客户端
func (m *Manager) initializeClients() error {
	if m.config.Databases == nil {
		return fmt.Errorf("no database configurations found")
	}

	var clientNames []string
	for name, dbConfig := range m.config.Databases {
		client, err := m.createClient(name, dbConfig)
		if err != nil {
			return fmt.Errorf("failed to create database client for '%s': %w", name, err)
		}

		m.clients.Store(name, client)
		clientNames = append(clientNames, name)
	}

	m.logger.Info("Database manager initialized successfully",
		zap.Int("client_count", len(clientNames)),
		zap.Strings("clients", clientNames),
		zap.String("default_database", m.defaultName),
		zap.Int("alias_count", len(m.aliases)))

	return nil
}

// createClient 创建数据库客户端（合并Builder功能）
func (m *Manager) createClient(name string, cfg config.DatabaseConfig) (*Client, error) {
	// 生成DSN
	dsn, dbType, err := m.buildDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DSN: %w", err)
	}

	// 打开数据库连接
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池
	m.configureConnectionPool(db, cfg)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建 Ent 驱动
	driver := entsql.OpenDB(dbType, db)

	client := &Client{
		name:   name,
		driver: driver,
		db:     db,
		config: cfg,
		logger: m.logger,
	}

	m.logger.Info("Database client created successfully",
		zap.String("name", name),
		zap.String("type", cfg.Type),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database))

	return client, nil
}

// buildDSN 生成数据库连接字符串（合并ConfigManager功能）
func (m *Manager) buildDSN(cfg config.DatabaseConfig) (string, string, error) {
	switch cfg.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		return dsn, dialect.MySQL, nil
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
		return dsn, dialect.Postgres, nil
	case "sqlite":
		return cfg.Database, dialect.SQLite, nil
	default:
		return "", "", fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}

// configureConnectionPool 配置连接池（合并Builder功能）
func (m *Manager) configureConnectionPool(db *sql.DB, cfg config.DatabaseConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
}

// GetClient 获取指定名称的数据库客户端
func (m *Manager) GetClient(name string) (*Client, error) {
	if client, ok := m.clients.Load(name); ok {
		return client.(*Client), nil
	}
	return nil, fmt.Errorf("database client '%s' not found", name)
}

// GetByAlias 通过别名获取数据库客户端
func (m *Manager) GetByAlias(alias string) (*Client, error) {
	if dbName, exists := m.aliases[alias]; exists {
		return m.GetClient(dbName)
	}
	return nil, fmt.Errorf("database alias '%s' not found", alias)
}

// Default 获取默认数据库客户端
func (m *Manager) Default() (*Client, error) {
	if m.defaultName == "" {
		return nil, fmt.Errorf("no default database configured")
	}
	return m.GetClient(m.defaultName)
}

// ListClients 列出所有数据库客户端名称
func (m *Manager) ListClients() []string {
	var names []string
	m.clients.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

// ListAliases 列出所有别名映射
func (m *Manager) ListAliases() map[string]string {
	result := make(map[string]string)
	for alias, dbName := range m.aliases {
		result[alias] = dbName
	}
	return result
}

// HasClient 检查是否存在指定名称的客户端
func (m *Manager) HasClient(name string) bool {
	_, ok := m.clients.Load(name)
	return ok
}

// Close 关闭所有数据库连接
func (m *Manager) Close() error {
	var lastErr error
	m.clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		if err := client.Close(); err != nil {
			m.logger.Error("Failed to close database client",
				zap.String("name", key.(string)),
				zap.Error(err))
			lastErr = err
		}
		return true
	})
	return lastErr
}
