package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// RedisLoginRateLimiter enforces per-IP login attempt limits using Redis.
type RedisLoginRateLimiter struct {
	client *redis.Client
	max    int64
	window time.Duration
	prefix string
}

func NewRedisLoginRateLimiter(redisURL string, maxAttempts int, window time.Duration) (*RedisLoginRateLimiter, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if pingErr := client.Ping(ctx).Err(); pingErr != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", pingErr)
	}

	return &RedisLoginRateLimiter{
		client: client,
		max:    int64(maxAttempts),
		window: window,
		prefix: "fidely:ratelimit:login:ip",
	}, nil
}

func (limiter *RedisLoginRateLimiter) Close() error {
	return limiter.client.Close()
}

// Ping verifies Redis connectivity for operational health checks.
func (limiter *RedisLoginRateLimiter) Ping(ctx context.Context) error {
	return limiter.client.Ping(ctx).Err()
}

func (limiter *RedisLoginRateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := fmt.Sprintf("%s:%s", limiter.prefix, c.RealIP())
			count, ttl, err := limiter.bump(c.Request().Context(), key)
			if err != nil {
				return c.JSON(http.StatusServiceUnavailable, map[string]any{
					"message": "login unavailable right now",
				})
			}
			if count > limiter.max {
				retryAfter := int64(ttl.Seconds())
				if retryAfter < 1 {
					retryAfter = 1
				}
				c.Response().Header().Set("Retry-After", strconv.FormatInt(retryAfter, 10))
				return c.JSON(http.StatusTooManyRequests, map[string]any{
					"message": "too many login attempts, please try again later",
				})
			}
			return next(c)
		}
	}
}

func (limiter *RedisLoginRateLimiter) bump(ctx context.Context, key string) (int64, time.Duration, error) {
	pipeline := limiter.client.TxPipeline()
	incr := pipeline.Incr(ctx, key)
	ttlCmd := pipeline.TTL(ctx, key)
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return 0, 0, err
	}

	count := incr.Val()
	ttl := ttlCmd.Val()
	if count == 1 || ttl <= 0 {
		if expireErr := limiter.client.Expire(ctx, key, limiter.window).Err(); expireErr != nil {
			return 0, 0, expireErr
		}
		ttl = limiter.window
	}

	return count, ttl, nil
}
