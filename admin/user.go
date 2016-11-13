package admin

import (
	"net/http"

	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/util"
)

// CreateUserHandler allows for the creation of users
func (a *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["test"] = "testing"

	w.Write(util.JSONify(root))
}

// ReplaceUserHandler allows for the modification of users
func (a *Handler) ReplaceUserHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["test"] = "testing"

	w.Write(util.JSONify(root))
}
