package contextutil

import (
	"context"
)

// 定义不导出的自定义类型作为context的key，避免键冲突
type contextKey string

// 用于在context中共享数据的键
var (
	// UserIDKey 是在context中存储用户ID的键
	UserIDKey contextKey = "user_id"
	// ClientIPContextKey 是在context中存储客户端IP字符串的键
	ClientIPContextKey contextKey = "clientIP"
	// ClientParsedIPContextKey 是在context中存储解析后的net.IP对象的键
	ClientParsedIPContextKey contextKey = "clientParsedIP"
	// TraceIDKey 是在context中存储追踪ID的键
	TraceIDKey contextKey = "traceID"
	// AuthHeaderKey 认证头键名
	AuthHeaderKey contextKey = "Authorization"
	// ContextLoggerKey 存储日志记录器的键
	ContextLoggerKey contextKey = "contextLogger"
)

// TokenPrefix Token前缀
const TokenPrefix = "Bearer "

// GetUserIDFromContext 从context中获取用户ID
// 返回用户ID以及一个布尔值，表示是否成功找到了用户ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
