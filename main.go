package main

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/mattgen88/blog/handlers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

const (
	postTable = ``
	userTable = ``
	categoryTable = ``
)

func main() {
	// Set configuration file information
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/blog/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config /etc/blog/config: %s", err)
	}

	// Gather configuration
	dbFile := viper.GetString("dbfile")

	viper.SetDefault("port", 8088)
	port := viper.GetInt("port")

	viper.SetDefault("host", "127.0.0.1")
	host := viper.GetString("host")

	viper.SetDefault("initialize", false)
	viper.BindEnv("initialize")
	init := viper.GetBool("initialize")


	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// If in initilization mode, run initializeBlog
	if (init) {
		initializeBlog(db)
		return
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	h := handlers.New(r, db)

	r.HandleFunc("/articles/{category}", h.CategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[a-zA-Z-_]+}", h.ArticleHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d",host, port), r))
}

// Called to initialize the database
func initializeBlog(db *sql.DB) {
	fmt.Println(db)
	fmt.Println(postTable)
	fmt.Println(userTable)
	fmt.Println(categoryTable)
}
