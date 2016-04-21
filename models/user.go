package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User interface {
	SetPassword(string)
	SetRealName(string)
	Exists() bool
	Authenticate(string) bool
	IsAuthenticated() bool
	HasRole(string)
	Populate() error
	Save()
}

// SQL based User model
type SqlUser struct {
	db            *sql.DB
	Id            int
	Username      string
	pwhash        string
	Realname      string
	Role          string
	Created       time.Time
	Email         string
	authenticated bool
	dirty         bool
	populated     bool
}

// Create a User model
func NewSqlUser(username string, db *sql.DB) *SqlUser {
	u := &SqlUser{
		db:       db,
		Username: username,
	}
	if u.Exists() {
		u.Populate()
	}
	return u
}

// Set the password of the user
func (u *SqlUser) SetPassword(pw string) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("An error occurred encrypting")
	}
	u.pwhash = string(bs)
	u.dirty = true
}

// Set the real name of the user
func (u *SqlUser) SetRealName(rn string) {
	u.Realname = rn
	u.dirty = true
}

// Check if the user exists
func (u *SqlUser) Exists() bool {
	var count int
	err := u.db.QueryRow(`SELECT COUNT(*) FROM Users WHERE Username = ?`, u.Username).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			fmt.Println("Something went wrong")
		}
	}
	return true
}

// Authenticate the user
func (u *SqlUser) Authenticate(pw string) bool {
	if !u.populated {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.pwhash), []byte(pw))
	return err == nil
}

// Check if user is authenticated
func (u *SqlUser) IsAuthenticated() bool {
	return u.authenticated
}

// Set the user role
func (u *SqlUser) SetRole(role string) {
	u.Role = role
}

// Check if user has a certain role
func (u *SqlUser) HasRole(role string) bool {
	return u.Role == role
}

// Fetch data and populate struct
func (u *SqlUser) Populate() error {
	if u.populated {
		// Don't repopulate a populated model
		return errors.New("Model already populated")
	}
	if u.dirty {
		// Don't populate a dirty model
		return errors.New("Model dirty")
	}
	// Fetch data and populate
	err := u.db.QueryRow(`SELECT UserId
	FROM Users
	WHERE Username = "mgeneral"`, u.Username).Scan(u.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return errors.New("User does not exist")
		}
		return errors.New("Unknown error occurred")
	}
	u.populated = true
	return nil
}

// Save struct back to database
func (u *SqlUser) Save() error{
	var result sql.Result
	var err error
	if u.populated {
		// update
		result, err = u.db.Exec(`UPDATE Users
		SET
			Hash = ?,
			RealName = ?,
			Email = ?,
			Role = (SELECT RoleID FROM Role where Name = ?)
		WHERE UserId = ?`, u.pwhash, u.Realname, u.Email, u.Role, u.Id)
	} else {
		created := time.Now().Format("2016-04-20 20:34:30")
		result, err = u.db.Exec(`INSERT INTO Users (
			Username,
			Hash,
			Created,
			RealName,
			Email,
			Role
		) VALUES (?, ?, ?, ?, ?, (SELECT RoleID FROM Role WHERE Name = ?))`, u.Username, u.pwhash, created, u.Realname, u.Email, u.Role)
	}
	if err != nil {
		// Some kind of failure
		fmt.Println(err)
		return errors.New("Unable to save user")
	}

	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		fmt.Println(err)
		return errors.New("Unable to verify save")
	}

	u.dirty = false
	return nil
}
