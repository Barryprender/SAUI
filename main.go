package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saui/handlers"
	"saui/middleware"
	"saui/statestore"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config{
		addr:   env("SAUI_ADDR", ":8080"),
		dbPath: env("SAUI_DB", "./saui.db"),
	}

	store, err := statestore.New(cfg.dbPath, logger)
	if err != nil {
		logger.Error("failed to open state store", "err", err)
		os.Exit(1)
	}
	defer store.Close()

	gw := statestore.NewGateway(store, logger)
	h := handlers.New(gw, logger)

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	session := middleware.Session
	page := func(hf http.HandlerFunc) http.Handler {
		return middleware.Chain(hf, session, middleware.CSRF)
	}

	mux.Handle("GET /{$}", page(h.Home))
	mux.Handle("GET /why", page(h.Why))
	mux.Handle("GET /architecture", page(h.Architecture))
	mux.Handle("GET /stack", page(h.Stack))
	mux.Handle("GET /cases", page(h.Cases))
	mux.Handle("GET /cases/food-ordering", page(h.CaseFoodOrdering))
	mux.Handle("GET /cases/banking", page(h.CaseBanking))
	mux.Handle("GET /cases/healthcare", page(h.CaseHealthcare))
	mux.Handle("GET /cases/saas-dashboard", page(h.CaseSaaS))
	mux.Handle("GET /cases/distributed-systems", page(h.CaseDistributed))
	mux.Handle("GET /cases/micro-frontends", page(h.CaseMFE))
	mux.Handle("GET /testing", page(h.Testing))
	mux.Handle("GET /limits", page(h.Limits))
	mux.Handle("GET /blog", page(h.Blog))
	mux.Handle("GET /code", page(h.Code))

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      middleware.Chain(mux, middleware.RateLimit, middleware.RecoverPanic(logger)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server starting", "addr", cfg.addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		logger.Error("shutdown error", "err", err)
	}
}

type config struct {
	addr   string
	dbPath string
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
