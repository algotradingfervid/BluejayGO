package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	visitors sync.Map
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{limit: limit, window: window}
	// Cleanup old entries periodically
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.visitors.Range(func(key, value interface{}) bool {
				v := value.(*visitor)
				if time.Since(v.lastSeen) > rl.window {
					rl.visitors.Delete(key)
				}
				return true
			})
		}
	}()
	return rl
}

func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			val, _ := rl.visitors.LoadOrStore(ip, &visitor{})
			v := val.(*visitor)

			now := time.Now()
			if now.Sub(v.lastSeen) > rl.window {
				v.count = 0
			}
			v.count++
			v.lastSeen = now

			if v.count > rl.limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests. Please try again later.",
				})
			}

			return next(c)
		}
	}
}
