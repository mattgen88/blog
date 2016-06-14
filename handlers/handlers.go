package handlers

import (
	"database/sql"
	"fmt"

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
