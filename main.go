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

	"github.com/mattgen88/blog/admin"
	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/internal/setup"
	"github.com/mattgen88/blog/util"
)

func main() {
	// Set configuration file information
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/blog/")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config /etc/blog/config: %s", err)
	}

	// Gather configuration
	viper.BindEnv("dbfile")
	dbFile := viper.GetString("dbfile")

	viper.SetDefault("port", 8088)
	port := viper.GetInt("port")

	viper.SetDefault("host", "127.0.0.1")
	host := viper.GetString("host")

	viper.SetDefault("initialize", false)
	viper.BindEnv("initialize")
	init := viper.GetBool("initialize")

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

	// If in initilization mode, run initializeBlog
	if init {
		setup.InitializeBlog(db)
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

	// Set up administrative access
	go func() {
		admin.Start(db)
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), util.ContentType(Gorilla.LoggingHandler(os.Stdout, r), "application/hal+json")))
}
