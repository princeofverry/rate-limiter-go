package httpapi

import (
	"encoding/json"
	"net/http"
	"princeofverry-rate-limiter/internal/apikey"
	"princeofverry-rate-limiter/internal/ratelimit"
)

type Middleware struct {
	KeyStore *apikey.Store
	Limiter *ratelimit.Limiter
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (m *Middleware) AuthAndLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]any {
				"error": "missing X-API-KEY",
			})
			return
		}
		if !m.KeyStore.Exists(key) {
			writeJSON(w, http.StatusUnauthorized, map[string]any {
				"error": "invalid api key",
			})
			return
		}
		if !m.Limiter.Allow(key) {
			writeJSON(w, http.StatusTooManyRequests, map[string]any {
				"error": "rate limit exceeded",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}