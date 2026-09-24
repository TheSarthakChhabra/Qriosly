package handler

import (
	"net"
	"net/http"
	"quiz-backend/apperror"
	"sync"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter *rate.Limiter
}

var (
	visitors   = make(map[string]*visitor)
	visitorsMu sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()
	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(12), 5)
		visitors[ip] = &visitor{limiter: limiter}
		return limiter
	}
	return v.limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			writeError(w, apperror.New(http.StatusTooManyRequests, "RATE_LIMITED", "too many attempts, please try again later"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
