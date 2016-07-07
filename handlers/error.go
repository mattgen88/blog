package handlers

import (
	"net/http"

	"github.com/mattgen88/blog/util"
	"github.com/pmoule/go2hal/hal"
)

// ErrorHandler handles requests for users
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["message"] = "Resource not found"

	w.WriteHeader(http.StatusNotFound)

	w.Write(util.JSONify(root))
}
