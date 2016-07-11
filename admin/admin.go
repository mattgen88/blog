package admin

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/util"
)

// Handler provides various http handlers
type AdminHandler struct {
	r  *mux.Router
	db *sql.DB
}

func New(r *mux.Router, db *sql.DB) *AdminHandler {
	return &AdminHandler{r, db}
}

// ArticleHandler should take posts of articles and save them to the database
// after checking for possible problems
func (a *AdminHandler) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["test"] = "testing"

	w.Write(util.JSONify(root))
}

// CategoryHandler should take posts of categories and save them to the database
// after checking for possible problems
func (a *AdminHandler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	w.Write(util.JSONify(root))
}

// CategoryHandler should take posts of users and save them to the database
// after checking for possible problems
func (a *AdminHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	w.Write(util.JSONify(root))
}

func Start(db *sql.DB) {
	r := mux.NewRouter()
	h := New(r, db)

	r.HandleFunc("/article", h.ArticleHandler)
	r.HandleFunc("/article/", h.ArticleHandler)
	r.HandleFunc("/user", h.UserHandler)
	r.HandleFunc("/user/", h.UserHandler)
	r.HandleFunc("/category", h.CategoryHandler)
	r.HandleFunc("/category/", h.CategoryHandler)

	r.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)
	// Firewall prevents access to this outside the network
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", Gorilla.LoggingHandler(os.Stdout, r)))
}
