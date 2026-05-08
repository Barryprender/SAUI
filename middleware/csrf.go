package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const csrfCookie = "saui_csrf"
const csrfHeader = "X-CSRF-Token"
const csrfField = "_csrf"

// CSRF enforces double-submit cookie pattern on state-mutating requests.
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := csrfToken(w, r)
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			submitted := r.Header.Get(csrfHeader)
			if submitted == "" {
				submitted = r.FormValue(csrfField)
			}
			if submitted == "" || submitted != token {
				http.Error(w, "invalid csrf token", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// CSRFToken returns the current CSRF token from the request context (set by csrfToken).
func CSRFToken(r *http.Request) string {
	if c, err := r.Cookie(csrfCookie); err == nil {
		return c.Value
	}
	return ""
}

func csrfToken(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(csrfCookie); err == nil {
		return c.Value
	}
	token := newToken()
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // must be readable by JS for header submission
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
	})
	return token
}

func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
