package middlewares

import (
	"net/http"
	"strings"
)

type RateLimiter struct {
	redisClient *redis.Client
	limiter     *rate.Limiter
}

func RateLimiter(resid *redis.Client) *RateLimiter {
	return &RateLimiter{
		redisClient: redis,
		limiter:     rate.RateLimiter(10, 5),
	}
}

func (rl *RateLimiter) LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr // Fallback to RemoteAddr if X-Forwarded-For is not present
		} else {
			// X-Forwarded-For can contain multiple IPs if multiple proxies are involved.
			// The first IP is typically the client's original IP.
			parts := strings.Split(clientIP, ",")
			clientIP = strings.TrimSpace(parts[0])
		}
		key := "rate_limit:" + clientIP
		count, err := rl.redisClient.Incr(r.Context(), key).Result()
		if err == nil && count > 10 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		rl.redisClient.Expire(r.Context(), key, time.Minute)

		
		if !rl.limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
