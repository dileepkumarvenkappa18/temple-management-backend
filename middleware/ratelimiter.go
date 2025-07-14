package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiterMiddleware returns a Gin middleware that limits requests per IP
func RateLimiterMiddleware() gin.HandlerFunc {
	// 🕒 Define the rate limit: 5 requests per second
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  5,
	}

	// 🧠 In-memory store (use Redis in prod if you need distributed limiting)
	store := memory.NewStore()

	// 📊 Limiter instance
	instance := limiter.New(store, rate)

	// 🚦 Gin-compatible middleware
	return ginlimiter.NewMiddleware(instance)
}
