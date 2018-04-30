package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// RootHandler handles requests for the root of the API
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self("/")
	root.AddLink("Users", &haljson.Link{Href: "/users"})
	root.AddLink("Article", &haljson.Link{Href: "/articles/{category}/{id:[a-zA-Z-_]+}", Templated: true})
	root.AddLink("Articles", &haljson.Link{Href: "/articles"})
	root.AddLink("Article Category", &haljson.Link{Href: "/articles/{category}", Templated: true})
	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
