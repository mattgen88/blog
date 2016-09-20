package util

import (
	"log"
	"net/http"

	"github.com/pmoule/go2hal/hal"
)

// JSONify the resource
func JSONify(root hal.Resource) []byte {

	encoder := new(hal.Encoder)
	bytes, err := encoder.ToJSON(root)

	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}

// ContentType sets the ContentType header to type
func ContentType(next http.Handler, ctype string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ctype)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
