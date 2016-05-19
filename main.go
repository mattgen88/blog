package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/howeyc/gopass"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/mattgen88/blog/handlers"
	"github.com/mattgen88/blog/models"
)

const (
	tableClean          = `DROP TABLE IF EXISTS Posts; DROP TABLE IF EXISTS Category; DROP TABLE IF EXISTS Users; DROP TABLE IF EXISTS Role;`
	categoryTableCreate = `CREATE TABLE Category (
		CategoryId Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Name Text NOT NULL
	);`
	categoryTableInsert = `INSERT INTO Category (
		"Name"
	) VALUES (
		"Test"
	);`
	roleTableCreate = `CREATE TABLE Role (
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
	userTableCreate = `CREATE TABLE Users (
		UserId Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Username Text NOT NULL,
		Hash TEXT NOT NULL,
		Created Datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		RealName Text,
		Email Text,
		Role Integer NULL REFERENCES Role(RoleID)
	);`
	postTableCreate = `CREATE TABLE Posts (
		PostId Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Title Text NOT NULL,
		Author Integer NOT NULL REFERENCES Users(UserId),
		Body Text NOT NULL,
		Date Datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		Slug Text NOT NULL,
		Category Integer NOT NULL DEFAULT 1 REFERENCES Category(CategoryId)
	);`
	postTableInsert = `INSERT INTO Posts (Title, Author, Body, Date, Slug, Category) VALUES ("Test", 1, "This is a test", CURRENT_TIMESTAMP, "test", 1);`
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

	r.HandleFunc("/", h.RootHandler)
	r.HandleFunc("/articles", h.ArticleListHandler)
	r.HandleFunc("/articles/{category}", h.CategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[a-zA-Z-_]+}", h.ArticleHandler)
	r.HandleFunc("/users/", h.UsersListHandler)
	r.HandleFunc("/users/{id:[a-zA-Z]}", h.UserHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), Gorilla.LoggingHandler(os.Stdout, r)))
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

	_, err := db.Exec(tableClean)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(categoryTableCreate)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(categoryTableInsert)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(roleTableCreate)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(roleTableInsert)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(userTableCreate)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("What is the username you would like to use? ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	username = strings.TrimSpace(username)

	fmt.Print("What is the real name you would like to use? ")
	realname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	realname = strings.TrimSpace(realname)

	fmt.Print("What is the email you would like to use? ")
	email, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	email = strings.TrimSpace(email)

	fmt.Print("What is the password for your user? ")

	password, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Println(err)
		return
	}

	u := models.NewSQLUser(username, db)
	u.SetPassword(strings.TrimSpace(string(password)))
	u.SetRealName(realname)
	u.SetEmail(email)
	err = u.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.Query("SELECT UserId FROM Users WHERE email=?", email)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(postTableCreate)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(postTableInsert, u.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done.")

}
