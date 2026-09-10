package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type ctxKey string

const sessionCtxKey ctxKey = "session_id"
const sessionCookieName = "saui_session"
const sessionIDBytes = 16 // produces 32 hex chars
const sessionIDHexLen = sessionIDBytes * 2

// NewSession returns the session middleware.
// secure must be true when serving over HTTPS or behind a TLS-terminating proxy.
func NewSession(secure bool, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := ""
			if c, err := r.Cookie(sessionCookieName); err == nil {
				if isValidSessionID(c.Value) {
					sessionID = c.Value
				} else {
					logger.Warn("rejecting malformed session id", "remote", r.RemoteAddr)
				}
			}
			if sessionID == "" {
				id, err := newSessionID()
				if err != nil {
					logger.Error("session id generation failed", "err", err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				sessionID = id
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,
					Value:    sessionID,
					Path:     "/",
					HttpOnly: true,
					Secure:   secure,
					SameSite: http.SameSiteStrictMode,
					Expires:  time.Now().Add(24 * time.Hour),
				})
			}
			ctx := context.WithValue(r.Context(), sessionCtxKey, sessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SessionID extracts the session ID from context.
func SessionID(ctx context.Context) string {
	v, _ := ctx.Value(sessionCtxKey).(string)
	return v
}

func isValidSessionID(id string) bool {
	if len(id) != sessionIDHexLen {
		return false
	}
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func newSessionID() (string, error) {
	b := make([]byte, sessionIDBytes)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
