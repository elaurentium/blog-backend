package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
	"sync"
)

type RateLimiterMiddleware struct {
	redisClient *redis.Client
	ipLimiters  map[string]*rate.Limiter
	mu          sync.Mutex
}

func NewRateLimiterMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		var count int
		err := redisClient.Get(ctx, key).Scan(&count)
		if err != nil {
			if err != redis.Nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}

			// Initial request
			count = 0
		}

		// Define o limite de requisi es (ex: 10 requisi es por minuto)
		limit := 10
		if count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		// Incrementa o contador no Redis
		err = redisClient.Set(ctx, key, count+1, time.Minute).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.Next()
	}
}

// Limiter para requisições por IP (sem Redis)
func (m *RateLimiterMiddleware) LimitByIP(rps int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		m.mu.Lock()
		limiter, exists := m.ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			m.ipLimiters[ip] = limiter
		}
		m.mu.Unlock()
		
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		
		c.Next()
	}
}

// Limiter distribuído com Redis
func (m *RateLimiterMiddleware) RedisRateLimiter(endpoint string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "ratelimit:" + endpoint + ":" + ip
		
		ctx := c.Request.Context()
		
		// Incrementa o contador no Redis
		val, err := m.redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next() // Em caso de erro, permitir a requisição
			return
		}
		
		// Se for a primeira requisição, defina o TTL
		if val == 1 {
			m.redisClient.Expire(ctx, key, window)
		}
		
		// Verificar se excedeu o limite
		if val > int64(limit) {
			ttl, _ := m.redisClient.TTL(ctx, key).Result()
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", ttl.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		
		c.Header("X-RateLimit-Remaining", string(int64(limit)-val))
		c.Next()
	}
}
