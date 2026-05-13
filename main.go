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
		addr:         env("SAUI_ADDR", ":8080"),
		dbPath:       env("SAUI_DB", "./saui.db"),
		secure:       env("SAUI_SECURE", "") == "true",
		trustedProxy: env("SAUI_TRUSTED_PROXY", ""),
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

	mux.Handle("GET /static/", http.StripPrefix("/static/",
		http.FileServer(noDirFS{http.Dir("static")})))

	session := middleware.NewSession(cfg.secure, logger)
	csrf := middleware.NewCSRF(cfg.secure, logger)
	page := func(hf http.HandlerFunc) http.Handler {
		return middleware.Chain(hf, session, csrf)
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
		Addr: cfg.addr,
		Handler: middleware.Chain(mux,
			middleware.SecurityHeaders(cfg.secure),
			middleware.NewRateLimit(cfg.trustedProxy, logger),
			middleware.RecoverPanic(logger),
		),
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
	addr         string
	dbPath       string
	secure       bool   // set via SAUI_SECURE=true when behind TLS
	trustedProxy string // set via SAUI_TRUSTED_PROXY=<ip> for X-Real-IP / X-Forwarded-For
}

// noDirFS wraps an http.FileSystem and returns ErrNotExist for directory requests,
// preventing directory listing of the static file tree.
type noDirFS struct{ http.FileSystem }

func (n noDirFS) Open(name string) (http.File, error) {
	f, err := n.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if s.IsDir() {
		f.Close()
		return nil, os.ErrNotExist
	}
	return f, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
