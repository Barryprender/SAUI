package middleware

import "net/http"

// HXRequest rejects requests to htmx-only routes that lack the HX-Request header.
func HXRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// IsHXRequest reports whether the request is an htmx partial request.
func IsHXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
