package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// RootHandler should take posts of articles and save them to the database
// after checking for possible problems
func (a *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self("/")
	root.AddLink("Auth", &haljson.Link{Href: "/auth"})
	root.AddLink("Refresh token", &haljson.Link{Href: "/refresh"})
	root.AddLink("Users", &haljson.Link{Href: "/users"})
	root.AddLink("User", &haljson.Link{Href: "/users/{user}", Templated: true})
	root.AddLink("Article", &haljson.Link{Href: "/articles/{id:[a-zA-Z-_]+}", Templated: true})
	root.AddLink("Articles", &haljson.Link{Href: "/articles"})
	root.AddLink("Articles by Category", &haljson.Link{Href: "/categories/{category}", Templated: true})
	root.AddLink("Categories", &haljson.Link{Href: "/categories"})
	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
