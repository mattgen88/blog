package admin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/util"
)

// Handler provides various http handlers
type Handler struct {
	r      *mux.Router
	db     *sql.DB
	jwtKey string
}

// New returns a new instance of the AdminHandler
func New(r *mux.Router, db *sql.DB, jwtKey string) *Handler {
	return &Handler{r, db, jwtKey}
}

// Start is called to configure and start the admin interface
func Start(db *sql.DB) {
	// Set configuration file information
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/blog/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config /etc/blog/config: %s", err)
	}

	viper.SetDefault("adminPort", 9001)
	adminPort := viper.GetInt("adminPort")

	viper.SetDefault("adminHost", "127.0.0.1")
	adminHost := viper.GetString("adminHost")

	jwtKey := viper.GetString("jwtKey")

	router := mux.NewRouter()
	h := New(router, db, jwtKey)
	ro := handlers.New(router, db)

	var articleHandlers Gorilla.MethodHandler
	articleHandlers = make(map[string]http.Handler)
	articleHandlers["GET"] = http.HandlerFunc(ro.ArticleHandler)
	articleHandlers["POST"] = AuthMiddleware(http.HandlerFunc(h.ReplaceArticleHandler), jwtKey, db)

	var articleListHandlers Gorilla.MethodHandler
	articleListHandlers = make(map[string]http.Handler)
	articleListHandlers["GET"] = http.HandlerFunc(ro.ArticleListHandler)
	articleListHandlers["POST"] = AuthMiddleware(http.HandlerFunc(h.CreateArticleHandler), jwtKey, db)

	var categoryListHandlers Gorilla.MethodHandler
	categoryListHandlers = make(map[string]http.Handler)
	categoryListHandlers["GET"] = http.HandlerFunc(ro.CategoryListHandler)
	categoryListHandlers["POST"] = AuthMiddleware(http.HandlerFunc(h.CreateCategoryHandler), jwtKey, db)

	var categoryHandlers Gorilla.MethodHandler
	categoryHandlers = make(map[string]http.Handler)
	categoryHandlers["GET"] = http.HandlerFunc(ro.CategoryHandler)
	categoryHandlers["POST"] = AuthMiddleware(http.HandlerFunc(h.ReplaceCategoryHandler), jwtKey, db)

	var userHandlers Gorilla.MethodHandler
	userHandlers = make(map[string]http.Handler)
	userHandlers["GET"] = http.HandlerFunc(ro.UserHandler)
	userHandlers["POST"] = http.HandlerFunc(h.ReplaceUserHandler)

	var userListHandlers Gorilla.MethodHandler
	userListHandlers = make(map[string]http.Handler)
	userListHandlers["GET"] = http.HandlerFunc(ro.UsersListHandler)
	userListHandlers["POST"] = http.HandlerFunc(h.CreateUserHandler)

	router.HandleFunc("/", h.RootHandler)

	router.Handle("/articles", articleListHandlers)
	router.Handle("/articles/", articleListHandlers)

	router.Handle("/categories", categoryListHandlers)
	router.Handle("/categories/", categoryListHandlers)

	router.Handle("/categories/{category}", categoryHandlers)
	router.Handle("/categories/{category}/", categoryHandlers)

	router.Handle("/articles/{id}", articleHandlers)
	router.Handle("/articles/{id}/", articleHandlers)

	router.Handle("/users", AuthMiddleware(userListHandlers, jwtKey, db))
	router.Handle("/users/", AuthMiddleware(userListHandlers, jwtKey, db))

	router.Handle("/users/{id}", AuthMiddleware(userHandlers, jwtKey, db))
	router.Handle("/users/{id}/", AuthMiddleware(userHandlers, jwtKey, db))

	router.Handle("/auth", http.HandlerFunc(h.Auth))
	router.Handle("/auth/", http.HandlerFunc(h.Auth))
	router.Handle("/refresh", AuthMiddleware(http.HandlerFunc(h.AuthRefresh), jwtKey, db))
	router.Handle("/refresh/", AuthMiddleware(http.HandlerFunc(h.AuthRefresh), jwtKey, db))
	router.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	// Firewall prevents access to this outside the network
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", adminHost, adminPort), util.ContentType(Gorilla.LoggingHandler(os.Stdout, router), "application/hal+json")))
}
