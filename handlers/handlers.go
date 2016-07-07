package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
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
