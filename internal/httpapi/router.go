package httpapi

import "net/http"

func NewRouter(h *Handlers, mw *Middleware) http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)

	// key management
	mux.HandleFunc("POST /v1/keys", h.CreateKey)
	// token bucket mutable dan concurrent 
	// jadi perlu membuat middleware auth only untuk handler LimitStatus
	mux.Handle("/v1/limit", mw.AuthOnly(http.HandlerFunc(h.LimitStatus)))
	mux.HandleFunc("DELETE /v1/keys/{key}", h.RevokeKey)

	// protected
	mux.Handle("/v1/ping", mw.AuthAndLimit(http.HandlerFunc(h.Ping)))

	return mux
}