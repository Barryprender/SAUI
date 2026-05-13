package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
)

const csrfCookieName = "saui_csrf"
const csrfHeader     = "X-CSRF-Token"
const csrfField      = "_csrf"

type csrfCtxKey struct{}

// NewCSRF returns the CSRF middleware using the double-submit cookie pattern.
// secure must be true when serving over HTTPS or behind a TLS-terminating proxy.
func NewCSRF(secure bool, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := ensureCSRFToken(w, r, secure)
			if err != nil {
				logger.Error("csrf token generation failed", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), csrfCtxKey{}, token))

			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				submitted := r.Header.Get(csrfHeader)
				if submitted == "" {
					submitted = r.FormValue(csrfField)
				}
				if submitted == "" || submitted != token {
					logger.Warn("csrf validation failed",
						"method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
					http.Error(w, "invalid csrf token", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CSRFToken returns the CSRF token from the request context.
// Returns "" if NewCSRF middleware has not run for this request.
func CSRFToken(r *http.Request) string {
	v, _ := r.Context().Value(csrfCtxKey{}).(string)
	return v
}

func ensureCSRFToken(w http.ResponseWriter, r *http.Request, secure bool) (string, error) {
	if c, err := r.Cookie(csrfCookieName); err == nil && c.Value != "" {
		return c.Value, nil
	}
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // JS-readable for header-based CSRF submission
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
	return token, nil
}
