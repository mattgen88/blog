package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/howeyc/gopass"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/models"
)

const (
	categoryTableCreate = `DROP TABLE IF EXISTS Category; CREATE TABLE Category (
		CategoryId Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Name Text NOT NULL
	);`
	categoryTableInsert = `INSERT INTO Category (
		"Name"
	) VALUES (
		"Test"
	);`
	roleTableCreate = `DROP TABLE IF EXISTS Role; CREATE TABLE Role (
		RoleID Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Name TEXT NOT NULL
	);`
	roleTableInsert = `INSERT INTO Role (
		"Name"
	) VALUES (
		"admin"
	);
	INSERT INTO Role (
		"Name"
	) VALUES (
		"user"
	);`
	userTableCreate = `DROP TABLE IF EXISTS Users; CREATE TABLE Users (
		UserID Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Username Text NOT NULL,
		Hash TEXT NOT NULL,
		Created Datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		RealName Text,
		Email Text,
		Role Integer NULL REFERENCES Role(RoleID)
	);`
	postTableCreate = `DROP TABLE IF EXISTS Posts; CREATE TABLE Posts (
		PostId Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Title Text NOT NULL,
		Author Integer NOT NULL REFERENCES Users(UserId),
		Body Text NOT NULL,
		Date Datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		Slug Text NOT NULL,
		Category Integer NOT NULL DEFAULT 1 REFERENCES Category(CategoryId)
	);`
	postTableInsert = `INSERT INTO Posts (
		Title,
		Author,
		Body,
		Date,
		Slug
		Category
	) VALUES (
		"Test",
		?,
		"This is a test",
		CURRENT_TIMESTAMP,
		"test",
		1
	);`
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

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		fmt.Println("Error enabling foreign keys")
		return
	}

	// If in initilization mode, run initializeBlog
	if init {
		initializeBlog(db)
		return
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	h := handlers.New(r, db)

	r.HandleFunc("/articles/{category}", h.CategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[a-zA-Z-_]+}", h.ArticleHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r))
}

// Called to initialize the database
func initializeBlog(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`Welcome to blog initlialization. If you have already run this
please remove the INITIALIZE environment variable of set it to FALSE and restart.

If this is your first time running, this set up will prompt you for some
information in order to set up the blog for the first time. Follow the
instructions and you will have working blog.`)

	fmt.Println("Initializing tables...")

	_, err := db.Exec(categoryTableCreate)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not create category table")
		return
	}

	_, err = db.Exec(categoryTableInsert)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not populate category table")
		return
	}

	_, err = db.Exec(roleTableCreate)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not create roles table")
		return
	}

	_, err = db.Exec(roleTableInsert)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not populate roles table")
		return
	}

	_, err = db.Exec(userTableCreate)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not create user table")
		return
	}

	fmt.Print("What is the username you would like to use? ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		fmt.Println("Something has gone wrong. Exiting")
		return
	}
	username = strings.TrimSpace(username)

	fmt.Print("What is the password for your user? ")

	password, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Something has gone wrong. Exiting")
		return
	}

	u := models.NewSqlUser(username, db)
	u.SetPassword(strings.TrimSpace(string(password)))
	err = u.Save()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not save user")
		return
	}
	u = models.NewSqlUser(username, db)
	err = u.Populate()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not populate database")
		return
	}


	_, err = db.Exec(postTableCreate)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not create post table")
		return
	}
	_, err = db.Exec(postTableInsert, u.Id)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not create post table")
		return
	}

	fmt.Println("Done.")

}
