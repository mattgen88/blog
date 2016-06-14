package handlers

import (
	"fmt"
	"net/http"

	"github.com/pmoule/go2hal/hal"
)

// RootHandler handles requests for the root of the API
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	self := hal.NewSelfLinkRelation()
	self.SetLink(&hal.LinkObject{Href: r.URL.Path})

	root.AddLink(self)

	// Users
	users, err := hal.NewLinkRelation("Users")
	if err != nil {
		// Error creating link relation
		fmt.Println(err)
		return
	}
	users.SetLink(&hal.LinkObject{Href: "/users/"})

	root.AddLink(users)

	// Article
	article, err := hal.NewLinkRelation("Article")
	if err != nil {
		// Error creating link relation
		fmt.Println(err)
		return
	}
	article.SetLink(&hal.LinkObject{Href: "/articles/{category}/{id:[a-zA-Z-_]+}"})

	root.AddLink(article)

	// Article Category
	category, err := hal.NewLinkRelation("Article Category")
	if err != nil {
		// Error creating link relation
		fmt.Println(err)
		return
	}
	category.SetLink(&hal.LinkObject{Href: "/articles/{category}"})

	root.AddLink(category)

	// Write it out
	w.Write(JSONify(root))
}
