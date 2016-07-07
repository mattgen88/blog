package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/mattgen88/blog/admin"
	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/middleware"
	"github.com/mattgen88/blog/setup"
)

func main() {
	// Set configuration file information
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/blog/")
	viper.AddConfigPath(".")

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

	r.HandleFunc("/articles/{category}", h.CategoryHandler)
	r.HandleFunc("/articles/{category}/", h.CategoryHandler)

	r.HandleFunc("/articles/{category}/{id:[a-zA-Z-_]+}", h.ArticleHandler)
	r.HandleFunc("/articles/{category}/{id:[a-zA-Z-_]+}/", h.ArticleHandler)

	r.HandleFunc("/users", h.UsersListHandler)
	r.HandleFunc("/users/", h.UsersListHandler)

	r.HandleFunc("/users/{id:[a-zA-Z0-9]+}", h.UserHandler)
	r.HandleFunc("/users/{id:[a-zA-Z0-9]+}/", h.UserHandler)

	r.Handle("/authtest", middleware.AuthHandler(http.HandlerFunc(h.AuthTest)))

	r.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	// Set up administrative access
	go func() {
		admin.Start(db)
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), Gorilla.LoggingHandler(os.Stdout, r)))
}
