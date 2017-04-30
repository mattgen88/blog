package middleware

import (
	"net/http"
)

func AuthHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Snag JWT, or redirect to auth endpoint
		handler.ServeHTTP(w, r)
	})
}
