package httpapi

import (
	"encoding/json"
	"net/http"
	"princeofverry-rate-limiter/internal/apikey"
	"princeofverry-rate-limiter/internal/ratelimit"
	"strings"
)

type Handlers struct {
	KeyStore *apikey.Store
	Limiter *ratelimit.Limiter
}

func (h *Handlers) CreateKey(w http.ResponseWriter, r *http.Request) {
	key, err := h.KeyStore.Create()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any {
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any {
		"api_key": key,
	})
}

func (h *Handlers) RevokeKey(w http.ResponseWriter, r *http.Request) {
	// path: /v1/keys/{key}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 3 {
		writeJSON(w, http.StatusBadRequest, map[string]any {
			"error": "invalid path",
		})
		return
	}
	key := parts[2]

	ok := h.KeyStore.Revoke(key)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any {
			"error": "key not found",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"revoked": true})
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"message": "pong"})
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (h *Handlers) LimitStatus(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-API-Key")
  
	st, ok := h.Limiter.Status(key)
	if !ok {
	  writeJSON(w, http.StatusNotFound, map[string]any{
		"error": "rate limit bucket not initialized yet",
	  })
	  return
	}
  
	writeJSON(w, http.StatusOK, st)
  }
  