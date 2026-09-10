package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saui/actions"
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

	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/robots.txt")
	})
	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/sitemap.xml")
	})

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
	mux.Handle("GET /blog/cra-clock", page(h.BlogPostCRAClock))
	mux.Handle("GET /blog/supply-chain", page(h.BlogPostSupplyChain))
	mux.Handle("GET /blog/server-response-time", page(h.BlogPostServerResponseTime))
	mux.Handle("GET /blog/eu-compliance", page(h.BlogPostEUCompliance))
	mux.Handle("GET /code", page(h.Code))

	mux.Handle("GET /es/{$}", page(h.HomeES))
	mux.Handle("GET /es/why", page(h.WhyES))
	mux.Handle("GET /es/architecture", page(h.ArchitectureES))
	mux.Handle("GET /es/stack", page(h.StackES))
	mux.Handle("GET /es/cases", page(h.CasesES))
	mux.Handle("GET /es/cases/food-ordering", page(h.CaseFoodOrderingES))
	mux.Handle("GET /es/cases/banking", page(h.CaseBankingES))
	mux.Handle("GET /es/cases/healthcare", page(h.CaseHealthcareES))
	mux.Handle("GET /es/cases/saas-dashboard", page(h.CaseSaaSES))
	mux.Handle("GET /es/cases/distributed-systems", page(h.CaseDistributedES))
	mux.Handle("GET /es/cases/micro-frontends", page(h.CaseMFEES))
	mux.Handle("GET /es/testing", page(h.TestingES))
	mux.Handle("GET /es/limits", page(h.LimitsES))
	mux.Handle("GET /es/blog", page(h.BlogES))
	mux.Handle("GET /es/blog/cra-clock", page(h.BlogPostCRAClockES))
	mux.Handle("GET /es/blog/supply-chain", page(h.BlogPostSupplyChainES))
	mux.Handle("GET /es/blog/server-response-time", page(h.BlogPostServerResponseTimeES))
	mux.Handle("GET /es/blog/eu-compliance", page(h.BlogPostEUComplianceES))
	mux.Handle("GET /es/code", page(h.CodeES))

	mux.Handle("POST /feedback", page(h.SubmitFeedback))
	mux.Handle("GET /architecture/step/{step}", http.HandlerFunc(h.ArchStep))

	mux.Handle("GET /demo/food-ordering", page(h.FoodMenu))
	mux.Handle("POST /demo/food-ordering/cart/add", page(h.FoodCartAdd))
	mux.Handle("POST /demo/food-ordering/cart/remove", page(h.FoodCartRemove))
	mux.Handle("POST /demo/food-ordering/checkout", page(h.FoodCheckout))
	mux.Handle("GET /demo/food-ordering/confirmation", page(h.FoodConfirmation))

	mux.Handle("GET /demo/banking", page(h.BankingDashboard))
	mux.Handle("GET /demo/banking/transfer", page(h.BankingTransfer))
	mux.Handle("POST /demo/banking/transfer", page(h.BankingTransferSubmit))
	mux.Handle("GET /demo/banking/confirmation", page(h.BankingConfirmation))

	mux.Handle("GET /demo/healthcare", page(h.HealthcareSchedule))
	mux.Handle("POST /demo/healthcare/book", page(h.HealthcareBook))
	mux.Handle("POST /demo/healthcare/cancel", page(h.HealthcareCancel))

	mux.Handle("GET /demo/saas", page(h.SAASDashboard))
	mux.Handle("POST /demo/saas/invite", page(h.SAASInvite))
	mux.Handle("POST /demo/saas/archive", page(h.SAASArchive))

	mux.Handle("GET /demo/distributed", page(h.DistributedPipeline))
	mux.Handle("POST /demo/distributed/advance", page(h.DistributedAdvance))

	mux.Handle("GET /demo/mfe", page(h.MFEWorkspace))
	mux.Handle("POST /demo/mfe/toggle", page(h.MFEToggle))
	mux.Handle("POST /demo/mfe/create", page(h.MFECreate))

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

	// Storage limitation is a schedule, not an intention. Feedback is retained
	// with its session identifier cleared; every other event past the window
	// goes. Run once at boot so a restarted instance is compliant immediately.
	purge := func() {
		pctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		anonymised, deleted, err := store.PurgeExpired(pctx, time.Now().UTC(), actions.FeedbackSubmittedType)
		if err != nil {
			logger.Error("event retention purge failed", "err", err)
			return
		}
		logger.Info("event retention purge", "anonymised", anonymised, "deleted", deleted)
	}
	go func() {
		purge()
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				purge()
			}
		}
	}()

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
