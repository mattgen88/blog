package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/blog/util"
)

// Handler provides various http handlers
type Handler struct {
	r  *mux.Router
	db *sql.DB
}

// New returns a new instance of the AdminHandler
func New(r *mux.Router, db *sql.DB) *Handler {
	return &Handler{r, db}
}

// RootHandler should take posts of articles and save them to the database
// after checking for possible problems
func (a *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	self := hal.NewSelfLinkRelation()
	self.SetLink(&hal.LinkObject{Href: r.URL.Path})

	root.AddLink(self)

	// Users
	users, err := hal.NewLinkRelation("Users")

	if err != nil {
		// Error creating link relation
		log.Println(err)
		return
	}

	users.SetLink(&hal.LinkObject{Href: "/users/"})

	root.AddLink(users)

	// Article
	article, err := hal.NewLinkRelation("Article")

	if err != nil {
		// Error creating link relation
		log.Println(err)
		return
	}

	article.SetLink(&hal.LinkObject{Href: "/articles/{category}/{id:[a-zA-Z-_]+}"})

	root.AddLink(article)

	// Article Category
	category, err := hal.NewLinkRelation("Article Category")

	if err != nil {
		// Error creating link relation
		log.Println(err)
		return
	}

	category.SetLink(&hal.LinkObject{Href: "/articles/{category}"})

	root.AddLink(category)

	// Write it out
	w.Write(util.JSONify(root))
}

// ArticleHandler should take posts of articles and save them to the database
// after checking for possible problems
func (a *Handler) ArticleHandler(w http.ResponseWriter, r *http.Request) {

	// Set up our hal resource
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	// Get the slug of the post we're dealing with
	slug := mux.Vars(r)["id"]

	// Fetch the requested article
	model := models.NewSQLArticle(slug, a.db)
	model.Populate()

	// Unpack posted data into model
	err := json.NewDecoder(r.Body).Decode(&model)

	if err != nil {
		// parse error
		root.Data()["error"] = ParseError
		w.Write(util.JSONify(root))
		return
	}

	err = model.Save()

	if err != nil {
		root.Data()["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(util.JSONify(root))
		return
	}

	// Write the model out
	root.Data()["title"] = model.Title
	root.Data()["author"] = model.Author
	root.Data()["body"] = model.Body
	root.Data()["slug"] = model.Slug
	root.Data()["date"] = model.Date
	root.Data()["id"] = model.ID
	root.Data()["category"] = model.Category

	w.Write(util.JSONify(root))
}

// CategoryHandler should take posts of categories and save them to the database
// after checking for possible problems
func (a *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["test"] = "testing"

	w.Write(util.JSONify(root))
}

// UserHandler allows for the modification of users
func (a *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["test"] = "testing"

	w.Write(util.JSONify(root))
}

// Start is called to configure and start the admin interface
func Start(db *sql.DB) {
	router := mux.NewRouter()
	h := New(router, db)
	ro := handlers.New(router, db)

	var articleHandlers Gorilla.MethodHandler
	articleHandlers = make(map[string]http.Handler)
	articleHandlers["GET"] = http.HandlerFunc(ro.ArticleHandler)
	articleHandlers["POST"] = http.HandlerFunc(h.ArticleHandler)

	var categoryHandlers Gorilla.MethodHandler
	categoryHandlers = make(map[string]http.Handler)
	categoryHandlers["GET"] = http.HandlerFunc(ro.CategoryHandler)
	categoryHandlers["POST"] = http.HandlerFunc(h.CategoryHandler)

	var userHandlers Gorilla.MethodHandler
	userHandlers = make(map[string]http.Handler)
	userHandlers["GET"] = http.HandlerFunc(ro.UserHandler)
	userHandlers["POST"] = http.HandlerFunc(h.UserHandler)

	router.HandleFunc("/", h.RootHandler)

	router.HandleFunc("/articles", ro.ArticleListHandler)
	router.HandleFunc("/articles/", ro.ArticleListHandler)

	router.Handle("/articles/{category}", categoryHandlers)
	router.Handle("/articles/{category}/", categoryHandlers)

	router.Handle("/articles/{category}/{id:[a-zA-Z-_]+}", articleHandlers)
	router.Handle("/articles/{category}/{id:[a-zA-Z-_]+}/", articleHandlers)

	router.HandleFunc("/users", ro.UsersListHandler)
	router.HandleFunc("/users/", ro.UsersListHandler)

	router.Handle("/users/{id:[a-zA-Z0-9]+}", userHandlers)
	router.Handle("/users/{id:[a-zA-Z0-9]+}/", userHandlers)
	//
	router.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	// Firewall prevents access to this outside the network
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", util.ContentType(Gorilla.LoggingHandler(os.Stdout, router), "application/hal+json")))
}
