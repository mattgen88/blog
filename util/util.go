package util

import "net/http"

// ContentType sets the ContentType header to type
func ContentType(next http.Handler, ctype string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ctype)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
