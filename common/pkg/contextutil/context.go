package contextutil

import (
	"context"
)

// 用于在context中共享数据的键
const (
	// UserIDKey 是在context中存储用户ID的键
	UserIDKey = "user_id"
	// ClientIPContextKey 是在context中存储客户端IP字符串的键
	ClientIPContextKey = "clientIP"
	// ClientParsedIPContextKey 是在context中存储解析后的net.IP对象的键
	ClientParsedIPContextKey = "clientParsedIP"
	// TraceIDKey 是在context中存储追踪ID的键
	TraceIDKey = "traceID"
	// AuthHeaderKey 认证头键名
	AuthHeaderKey = "Authorization"
	// TokenPrefix Token前缀
	TokenPrefix = "Bearer "
)

// GetUserIDFromContext 从context中获取用户ID
// 返回用户ID以及一个布尔值，表示是否成功找到了用户ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
