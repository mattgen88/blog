package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"

)

// Handler provides various http handlers
type Handler struct {
	r  *mux.Router
	db *sql.DB
}

// New returns a configured handler struct
func New(r *mux.Router, db *sql.DB) *Handler {
	return &Handler{r, db}
}

// JSONify the resource
func JSONify(root hal.Resource) []byte {

	encoder := new(hal.Encoder)
	bytes, err := encoder.ToJSON(root)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return bytes
}

// ErrorHandler handles requests for users
func (h *Handler) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["message"] = "Resource not found"

	w.WriteHeader(http.StatusNotFound)

	w.Write(JSONify(root))
}
