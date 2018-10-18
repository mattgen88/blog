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
	root.AddLink("User", &haljson.Link{Href: "/users/{user}", Templated: true})
	root.AddLink("Article", &haljson.Link{Href: "/articles/{id:[a-zA-Z-_]+}", Templated: true})
	root.AddLink("Articles", &haljson.Link{Href: "/articles"})
	root.AddLink("Article for Category", &haljson.Link{Href: "/categories/{category}", Templated: true})
	root.AddLink("Categories", &haljson.Link{Href: "/categories"})
	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
