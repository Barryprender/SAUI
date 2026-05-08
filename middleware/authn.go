package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

type ctxKey string

const sessionKey ctxKey = "session_id"

const cookieName = "saui_session"

// Session reads or creates a session cookie and puts the session ID in context.
func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := ""
		if c, err := r.Cookie(cookieName); err == nil {
			sessionID = c.Value
		}
		if sessionID == "" {
			sessionID = newSessionID()
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(24 * time.Hour),
			})
		}
		ctx := context.WithValue(r.Context(), sessionKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SessionID extracts the session ID from context.
func SessionID(ctx context.Context) string {
	v, _ := ctx.Value(sessionKey).(string)
	return v
}

func newSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
