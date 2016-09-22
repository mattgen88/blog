package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is a user in the system
type User interface {
	SetPassword(string) string
	SetRealName(string)
	SetEmail(string)
	Exists() bool
	Authenticate(string) bool
	IsAuthenticated() bool
	HasRole(string)
	Populate() error
	GetRealname() string
	GetRole() string
	GetEmail() string
	GetCreated() *time.Time
	GetUsername() string
	GetID() int
}

// SQLUser is a SQL based User model
type SQLUser struct {
	Db            *sql.DB    `json:"-"`
	ID            int        `json:"id,omitempty"`
	Username      string     `json:"username"`
	Realname      string     `json:"realname,omitempty"`
	Role          string     `json:"role,omitempty"`
	Created       *time.Time `json:"created,omitempty"`
	Email         string     `json:"email,omitempty"`
	pwhash        string
	authenticated bool
	dirty         bool
	populated     bool
}

// NewSQLUser Creates a User model
func NewSQLUser(username string, db *sql.DB) *SQLUser {
	u := &SQLUser{
		Db:       db,
		Username: username,
	}

	if u.Exists() {
		u.Populate()
	}

	return u
}

// GetRealname returns the real name of the user
func (u *SQLUser) GetRealname() string {
	if u.populated {
		return u.Realname
	}

	return ""
}

// GetRole returns the roles for the user
func (u *SQLUser) GetRole() string {
	if u.populated {
		return u.Role
	}

	return ""
}

// GetEmail returns the email
func (u *SQLUser) GetEmail() string {
	if u.populated {
		return u.Email
	}

	return ""
}

// GetCreated returns the time it was created
func (u *SQLUser) GetCreated() *time.Time {
	if u.populated {
		return u.Created
	}

	return nil
}

// GetUsername returns the username
func (u *SQLUser) GetUsername() string {
	if u.populated {
		return u.Username
	}

	return ""
}

// GetID returns the user's ID
func (u *SQLUser) GetID() int {
	if u.populated {
		return u.ID
	}

	return -1
}

// SetPassword sets the password of the user
func (u *SQLUser) SetPassword(pw string) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
	}

	u.pwhash = string(bs)
	u.dirty = true
	return string(bs)
}

// SetRealName sets the real name of the user
func (u *SQLUser) SetRealName(rn string) {
	u.Realname = rn
	u.dirty = true
}

// SetEmail sets the email of the user
func (u *SQLUser) SetEmail(email string) {
	u.Email = email
	u.dirty = true
}

// Exists Checks if the user exists
func (u *SQLUser) Exists() bool {
	var count int
	err := u.Db.QueryRow(`SELECT COUNT(*) FROM Users WHERE Username = ?`, u.Username).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}

	return true
}

// Authenticate authenticates the user
func (u *SQLUser) Authenticate(pw string) bool {
	if !u.populated {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.pwhash), []byte(pw))
	return err == nil
}

// IsAuthenticated checks if user is authenticated
func (u *SQLUser) IsAuthenticated() bool {
	return u.authenticated
}

// SetRole sets the user role
func (u *SQLUser) SetRole(role string) {
	u.Role = role
}

// HasRole checks if user has a certain role
func (u *SQLUser) HasRole(role string) bool {
	return u.Role == role
}

// Populate Fetches data and populates struct
func (u *SQLUser) Populate() error {
	if !u.Exists() {
		return errors.New("Instance does not exist")
	}

	if u.populated {
		// Don't repopulate a populated model
		return errors.New("Model already populated")
	}

	if u.dirty {
		// Don't populate a dirty model
		return errors.New("Model dirty")
	}

	// Fetch data and populate
	err := u.Db.QueryRow(`SELECT UserId, Created, RealName, Email, Role
	FROM Users
	WHERE Username = ?`, u.Username).Scan(&u.ID, &u.Created, &u.Realname, &u.Email, &u.Role)

	if err != nil {
		return errors.New("Unknown error occurred")
	}

	u.populated = true
	return nil
}
