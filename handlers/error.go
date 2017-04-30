package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// ErrorHandler handles requests for users
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["message"] = "Resource not found"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(json)
}
