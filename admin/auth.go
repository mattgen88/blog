package admin

import (
	"net/http"

	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/util"
)

// AuthHandler handles request to authenticate and will issue a JWT
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	w.WriteHeader(http.StatusNotFound)

	w.Write(util.JSONify(root))
}
