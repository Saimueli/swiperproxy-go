package middleware

import (
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

func RateLimit(next http.Handler, maxRequests int, window time.Duration) http.Handler {
	mu := sync.Mutex{}
	visitors := map[string]*visitor{}

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > window {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &visitor{}
			visitors[ip] = v
		}
		if time.Since(v.lastSeen) > window {
			v.count = 0
		}
		v.count++
		v.lastSeen = time.Now()
		current := v.count
		mu.Unlock()

		if current > maxRequests {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}