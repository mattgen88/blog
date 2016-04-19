package handlers

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/gorilla/mux"
)

type handler struct {
	r *mux.Router
	db *sql.DB
}

func New(r *mux.Router, db *sql.DB) *handler {
	return &handler{r, db}
}

func (h *handler) CategoryHandler(http.ResponseWriter, *http.Request) {
	fmt.Println("Category")
}

func (h *handler) ArticleHandler(http.ResponseWriter, *http.Request) {
	fmt.Println("Article")
}


