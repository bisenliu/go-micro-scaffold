package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"go.uber.org/zap"

	"common/config"
	"common/logger"
	"common/pkg/contextutil"
	"common/response"
)

// ipLimiter 存储每个IP的令牌桶和最后访问时间
type ipLimiter struct {
	bucket   *ratelimit.Bucket
	lastSeen time.Time
	mu       sync.Mutex
}

var (
	ipBuckets   sync.Map
	cleanupOnce sync.Once
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(cfg config.RateLimitConfig, baseLogger *zap.Logger) gin.HandlerFunc {
	// 如果未启用，返回一个空操作的中间件
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 确保清理goroutine只启动一次
	cleanupOnce.Do(func() {
		startCleanupGoroutine(cfg, baseLogger)
	})

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var ip string
		// 优先从 context 获取真实 IP
		if clientIP, exists := c.Get(contextutil.ClientIPContextKey); exists {
			ip, _ = clientIP.(string)
		}

		// 如果 context 中没有，则回退
		if ip == "" {
			logger.Warn(ctx, "Client IP not found in context, falling back to c.ClientIP() for rate limiting")
			ip = c.ClientIP()
		}

		if ip == "" {
			logger.Error(ctx, "Could not determine client IP for rate limiting")
			response.FailWithCode(c, response.CodeInternalError, "Could not determine client IP")
			c.Abort()
			return
		}

		var limiter *ipLimiter
		val, ok := ipBuckets.Load(ip)
		if ok {
			limiter = val.(*ipLimiter)
		} else {
			newLimiter := &ipLimiter{
				bucket: ratelimit.NewBucketWithQuantum(cfg.FillInterval, cfg.Capacity, cfg.Quantum),
			}
			val, _ = ipBuckets.LoadOrStore(ip, newLimiter)
			limiter = val.(*ipLimiter)
		}

		limiter.mu.Lock()
		limiter.lastSeen = time.Now()
		canTake := limiter.bucket.TakeAvailable(1) >= 1
		limiter.mu.Unlock()

		if !canTake {
			response.FailWithCode(c, response.CodeRateLimit, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

// startCleanupGoroutine 启动一个后台goroutine来定期清理过期的IP令牌桶
func startCleanupGoroutine(cfg config.RateLimitConfig, log *zap.Logger) {
	go func() {
		log.Info("Starting rate limit cleanup goroutine",
			zap.Duration("cleanup_interval", cfg.CleanupInterval),
			zap.Duration("bucket_expiry", cfg.BucketExpiry))

		ticker := time.NewTicker(cfg.CleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			var cleanedCount int
			ipBuckets.Range(func(key, value interface{}) bool {
				limiter := value.(*ipLimiter)
				limiter.mu.Lock()
				isExpired := time.Since(limiter.lastSeen) > cfg.BucketExpiry
				limiter.mu.Unlock()

				if isExpired {
					ipBuckets.Delete(key)
					cleanedCount++
				}
				return true // continue iteration
			})

			if cleanedCount > 0 {
				log.Info("Finished rate limit cleanup cycle", zap.Int("cleaned_buckets", cleanedCount))
			} else {
				log.Debug("Finished rate limit cleanup cycle, no expired buckets found")
			}
		}
	}()
}
