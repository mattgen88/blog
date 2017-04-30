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

	templated := true
	root.Self("/")
	root.AddLink("Users", &haljson.Link{Href: "/users"})
	root.AddLink("Article", &haljson.Link{Href: "/articles/{category}/{id:[a-zA-Z-_]+}", Templated: &templated})
	root.AddLink("Articles", &haljson.Link{Href: "/articles"})
	root.AddLink("Article Category", &haljson.Link{Href: "/articles/{category}", Templated: &templated})
	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
