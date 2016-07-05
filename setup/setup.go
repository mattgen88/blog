package setup

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mattgen88/blog/models"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tableClean          = `DROP TABLE IF EXISTS Articles; DROP TABLE IF EXISTS Category; DROP TABLE IF EXISTS Users; DROP TABLE IF EXISTS Role;`
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
		Created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		RealName Text,
		Email Text,
		Role Integer NULL REFERENCES Role(RoleID)
	);`
	userTableInsert = `INSERT INTO Users (
		Username,
		Hash,
		Created,
		RealName,
		Email,
		Role
	) VALUES (
		?,
		?,
		CURRENT_TIMESTAMP,
		?,
		?,
		(
			SELECT RoleID
			FROM Role
			WHERE Name = ?
		)
	);`

	articlesTableCreate = `CREATE TABLE Articles (
		ArticleId Integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		Title Text NOT NULL,
		Author Integer NOT NULL REFERENCES Users(UserId),
		Body Text NOT NULL,
		Date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		Slug Text NOT NULL,
		Category Integer NOT NULL DEFAULT 1 REFERENCES Category(CategoryId)
	);`
	articlesTableInsert = `INSERT INTO Articles (Title, Author, Body, Date, Slug, Category) VALUES ("Test", 1, "This is a test", CURRENT_TIMESTAMP, "test", 1);`
)

// Called to initialize the database
func InitializeBlog(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`Welcome to blog initlialization. If you have already run this
please remove the INITIALIZE environment variable of set it to FALSE and restart.

If this is your first time running, this set up will prompt you for some
information in order to set up the blog for the first time. Follow the
instructions and you will have working blog.`)

	fmt.Println("Initializing tables...")

	_, err := db.Exec(tableClean)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec(categoryTableCreate)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec(categoryTableInsert)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec(roleTableCreate)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec(roleTableInsert)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec(userTableCreate)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Print("What is the username you would like to use? ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	username = strings.TrimSpace(username)

	fmt.Print("What is the real name you would like to use? ")
	realname, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	realname = strings.TrimSpace(realname)

	fmt.Print("What is the email you would like to use? ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	email = strings.TrimSpace(email)

	fmt.Print("What is the password for your user? ")

	password, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Println(err)
		return
	}

	u := models.NewSQLUser(username, db)
	pwhash := u.SetPassword(strings.TrimSpace(string(password)))
	u.SetRealName(realname)
	u.SetEmail(email)

	_, err = db.Exec(userTableInsert, u.Username, pwhash, u.Realname, u.Email, u.Role)
	if err != nil {
		log.Println(err)
		return
	}

	rows, err := db.Query(`SELECT UserId
		FROM Users
		WHERE email=?`, email)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	_, err = db.Exec(articlesTableCreate)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Exec(articlesTableInsert, u.ID)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Done.")

}
