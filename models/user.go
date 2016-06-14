package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is a user in the system
type User interface {
	SetPassword(string)
	SetRealName(string)
	SetEmail(string)
	Exists() bool
	Authenticate(string) bool
	IsAuthenticated() bool
	HasRole(string)
	Populate() error
	Save()
	GetRealname() string
	GetRole() string
	GetEmail() string
	GetCreated() *time.Time
	GetUsername() string
	GetID() int
}

// SQLUser is a SQL based User model
type SQLUser struct {
	db            *sql.DB
	ID            int
	Username      string
	pwhash        string
	Realname      string
	Role          string
	Created       *time.Time
	Email         string
	authenticated bool
	dirty         bool
	populated     bool
}

// NewSQLUser Creates a User model
func NewSQLUser(username string, db *sql.DB) *SQLUser {
	u := &SQLUser{
		db:       db,
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
func (u *SQLUser) SetPassword(pw string) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("An error occurred encrypting")
	}

	u.pwhash = string(bs)
	u.dirty = true
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
	err := u.db.QueryRow(`SELECT COUNT(*) FROM Users WHERE Username = ?`, u.Username).Scan(&count)

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
	err := u.db.QueryRow(`SELECT UserId, Created, RealName, Email, Role
	FROM Users
	WHERE Username = ?`, u.Username).Scan(u.ID, u.Created, u.Realname, u.Email, u.Role)

	if err != nil {
		return errors.New("Unknown error occurred")
	}

	u.populated = true
	return nil
}

// Save saves struct back to database
func (u *SQLUser) Save() error {
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
		WHERE UserId = ?`, u.pwhash, u.Realname, u.Email, u.Role, u.ID)

	} else {

		result, err = u.db.Exec(`INSERT INTO Users (
			Username,
			Hash,
			Created,
			RealName,
			Email,
			Role
		) VALUES (?, ?, CURRENT_TIMESTAMP, ?, ?, (SELECT RoleID FROM Role WHERE Name = ?))`, u.Username, u.pwhash, u.Realname, u.Email, u.Role)
	}

	if err != nil {
		// Some kind of failure
		return errors.New("Unable to save user")
	}

	count, err := result.RowsAffected()

	if err != nil || count != 1 {
		return errors.New("Unable to verify save")
	}

	u.dirty = false
	return nil
}
