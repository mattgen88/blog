package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// AuthTest handles requests for users
func (h *Handler) AuthTest(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["message"] = "You should only see this after authenticating"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(json)
}
