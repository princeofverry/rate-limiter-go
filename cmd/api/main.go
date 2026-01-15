package main

import (
	"context"
	// "log"
	"net/http"
	"os"
	"os/signal"
	"princeofverry-rate-limiter/internal/apikey"
	"princeofverry-rate-limiter/internal/httpapi"
	"princeofverry-rate-limiter/internal/ratelimit"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	keyStore := apikey.NewStore()
	limiter := ratelimit.New(60, 60) // capacity 60, refill 60 per minute

	handlers := &httpapi.Handlers{
		KeyStore: keyStore, 
		Limiter: limiter,
	}
	mw := &httpapi.Middleware{KeyStore: keyStore, Limiter: limiter}
	
	router := httpapi.NewRouter(handlers, mw)

	handler := mw.Logger(router)
	srv := &http.Server{
		Addr: ":" + port,
		Handler: handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// create a context that listens for OS interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start server in a garoutine (so main can listen for shutdown signal)
	go func() {
		log.Info().Str("addr", srv.Addr).Msg("server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server failed to start")
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	// Give in-flight requests time to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Info().Msg("server stopped gracefully")
	}

}