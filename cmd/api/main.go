package main

import (
	"log"
	"net/http"
	"os"
	"princeofverry-rate-limiter/internal/apikey"
	"princeofverry-rate-limiter/internal/httpapi"
	"princeofverry-rate-limiter/internal/ratelimit"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	keyStore := apikey.NewStore()
	limiter := ratelimit.New(60, 60) // capacity 60, refill 60 per minute

	handlers := &httpapi.Handlers{KeyStore: keyStore, Limiter: limiter}
	mw := &httpapi.Middleware{KeyStore: keyStore, Limiter: limiter}
	
	router := httpapi.NewRouter(handlers, mw)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}