package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// registers with database/sql
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/util"
)

func main() {
	viper.AutomaticEnv()

	// Gather configuration
	viper.BindEnv("dsn")
	dsn := viper.GetString("dsn")

	viper.BindEnv("port")
	viper.SetDefault("port", "8088")
	port := viper.GetString("port")

	viper.BindEnv("host")
	viper.SetDefault("host", "127.0.0.1")
	host := viper.GetString("host")
	log.Println("Starting on ", host, " port ", port, " dsn ", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := mux.NewRouter()

	h := handlers.New(r, db)

	r.HandleFunc("/", h.RootHandler).Name("root")

	r.HandleFunc("/articles", h.ArticleListHandler)
	r.HandleFunc("/articles/", h.ArticleListHandler)

	r.HandleFunc("/categories", h.CategoryListHandler)
	r.HandleFunc("/categories/", h.CategoryListHandler)

	r.HandleFunc("/categories/{category}", h.CategoryHandler)
	r.HandleFunc("/categories/{category}/", h.CategoryHandler)

	r.HandleFunc("/articles/{id}", h.ArticleHandler)
	r.HandleFunc("/articles/{id}/", h.ArticleHandler)

	r.HandleFunc("/users", h.UsersListHandler)
	r.HandleFunc("/users/", h.UsersListHandler)

	r.HandleFunc("/users/{id}", h.UserHandler)
	r.HandleFunc("/users/{id}/", h.UserHandler)

	r.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	log.Fatal(http.ListenAndServe(net.JoinHostPort(host, port), util.ContentType(Gorilla.LoggingHandler(os.Stdout, Gorilla.CORS()(r)), "application/hal+json")))
}
