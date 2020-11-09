package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// registers with database/sql
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/util"
)

func main() {
	viper.AutomaticEnv()

	// Gather configuration
	viper.BindEnv("dbfile")
	dbFile := viper.GetString("dbfile")

	viper.SetDefault("port", 8088)
	port := viper.GetInt("port")

	viper.SetDefault("host", "127.0.0.1")
	host := viper.GetString("host")

	db, err := sql.Open("sqlite3", dbFile+"?parseTime=True")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Println("Error enabling foreign keys", err)
		return
	}

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

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), util.ContentType(Gorilla.LoggingHandler(os.Stdout, r), "application/hal+json")))
}
