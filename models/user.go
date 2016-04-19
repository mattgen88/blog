package models

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User interface {
	SetPassword(string)
	SetRealName(string)
	Exists() bool
	Authenticate(string) bool
	HasRole(string)
	Populate() error
	Save()
}

// SQL based User model
type SqlUser struct {
	db            *sql.DB
	id            int
	username      string
	pwhash        string
	realname      string
	role          int
	authenticated bool
	dirty         bool
	populated     bool
}

// Create a User model
func NewSqlUser(username string, db *sql.DB) *SqlUser {
	u := &SqlUser{
		db:       db,
		username: username,
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
	u.dirty = true
}

// Check if the user exists
func (u *SqlUser) Exists() bool {
	var count int
	err := u.db.QueryRow(`SELECT COUNT(*) FROM Users WHERE Username = ?`, u.username).Scan(&count)

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

// Check if user has a certain role
func (u *SqlUser) HasRole(role string) {
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
	u.populated = true
	return nil
}

// Save struct back to database
func (u *SqlUser) Save() {
	u.dirty = false
}
