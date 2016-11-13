package admin

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/mattgen88/blog/handlers"
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

// Start is called to configure and start the admin interface
func Start(db *sql.DB) {
	router := mux.NewRouter()
	h := New(router, db)
	ro := handlers.New(router, db)

	var articleHandlers Gorilla.MethodHandler
	articleHandlers = make(map[string]http.Handler)
	articleHandlers["GET"] = http.HandlerFunc(ro.ArticleHandler)
	articleHandlers["POST"] = http.HandlerFunc(h.ReplaceArticleHandler)

	var articleListHandlers Gorilla.MethodHandler
	articleListHandlers = make(map[string]http.Handler)
	articleListHandlers["GET"] = http.HandlerFunc(ro.ArticleListHandler)
	articleListHandlers["POST"] = http.HandlerFunc(h.CreateArticleHandler)

	var categoryListHandlers Gorilla.MethodHandler
	categoryListHandlers = make(map[string]http.Handler)
	categoryListHandlers["GET"] = http.HandlerFunc(ro.CategoryListHandler)
	categoryListHandlers["POST"] = http.HandlerFunc(h.CreateCategoryHandler)

	var categoryHandlers Gorilla.MethodHandler
	categoryHandlers = make(map[string]http.Handler)
	categoryHandlers["GET"] = http.HandlerFunc(ro.CategoryHandler)
	categoryHandlers["POST"] = http.HandlerFunc(h.ReplaceCategoryHandler)

	var userHandlers Gorilla.MethodHandler
	userHandlers = make(map[string]http.Handler)
	userHandlers["GET"] = http.HandlerFunc(ro.UserHandler)
	userHandlers["POST"] = http.HandlerFunc(h.ReplaceUserHandler)

	var userListHandlers Gorilla.MethodHandler
	userListHandlers = make(map[string]http.Handler)
	userListHandlers["GET"] = http.HandlerFunc(ro.UserHandler)
	userListHandlers["POST"] = http.HandlerFunc(h.CreateUserHandler)

	router.HandleFunc("/", h.RootHandler)

	router.Handle("/articles", articleListHandlers)
	router.Handle("/articles/", articleListHandlers)

	router.Handle("/categories", categoryListHandlers)
	router.Handle("/categories/", categoryListHandlers)

	router.Handle("/categories/{category}", categoryHandlers)
	router.Handle("/categories/{category}/", categoryHandlers)

	router.Handle("/articles/{id:[a-zA-Z-_]+}", articleHandlers)
	router.Handle("/articles/{id:[a-zA-Z-_]+}/", articleHandlers)

	router.Handle("/users", userListHandlers)
	router.Handle("/users/", userListHandlers)

	router.Handle("/users/{id:[a-zA-Z0-9]+}", userHandlers)
	router.Handle("/users/{id:[a-zA-Z0-9]+}/", userHandlers)
	//
	router.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	// Firewall prevents access to this outside the network
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", util.ContentType(Gorilla.LoggingHandler(os.Stdout, router), "application/hal+json")))
}
