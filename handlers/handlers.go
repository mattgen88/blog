package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

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

// CategoryHandler handles requests for categories
func (h *Handler) CategoryHandler(http.ResponseWriter, *http.Request) {
	fmt.Println("Category")
}

// ArticleHandler handles requests for articles
func (h *Handler) ArticleHandler(http.ResponseWriter, *http.Request) {
	fmt.Println("Article")
}

// UserHandler handles requests for users
func (h *Handler) UserHandler(http.ResponseWriter, *http.Request) {
	fmt.Println("User")
}
