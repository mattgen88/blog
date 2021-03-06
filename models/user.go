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
	Save() error
	Validate() error
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
	exists        bool
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
	if u.exists {
		return true
	}
	var count int
	err := u.Db.QueryRow(`SELECT COUNT(*) FROM "users" WHERE "username" = $1`, u.Username).Scan(&count)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	u.exists = count > 0
	return count > 0
}

// Authenticate authenticates the user
func (u *SQLUser) Authenticate(pw string) bool {
	if !u.populated {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.pwhash), []byte(pw))
	if err != nil {
		log.Println("mismatched hashing")
	}
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

	if u.dirty {
		// Don't populate a dirty model
		return errors.New("Model dirty")
	}

	// Fetch data and populate
	err := u.Db.QueryRow(`SELECT "userid", "created", "realname", "email", "role", "hash"
	FROM "users"
	WHERE "username" = $1`, u.Username).Scan(&u.ID, &u.Created, &u.Realname, &u.Email, &u.Role, &u.pwhash)

	if err != nil {
		log.Println(err)
		return errors.New("Unknown error occurred")
	}

	u.populated = true
	return nil
}

func (u *SQLUser) Save() error {
	var err error
	var query string

	err = u.Validate()
	if err != nil {
		// Validation error
		return err
	}
	if !u.Exists() {
		query = `INSERT INTO "users" (
			"username",
			"hash",
			"realName",
			"email",
			"created",
			"role"
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			CURRENT_TIMESTAMP,
			(
				SELECT "roleid"
				FROM "role"
				WHERE "name" = $5
			)
		);`
		result, err := u.Db.Exec(query, u.Username, u.pwhash, u.Realname, u.Email, "user")
		if err != nil {
			log.Println(err)
			return ErrSave
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Println(err)
			return ErrSave
		}
		u.ID = int(id)
	} else {
		query = `UPDATE "users" SET "hash" = $1, "realname" = $2, "email" = $3, "role" = $4 WHERE "userid" = $5`
		_, err = u.Db.Exec(query, u.pwhash, u.Realname, u.Email, u.Role, u.ID)
	}

	if err != nil {
		return ErrSave
	}

	return nil
}

func (u *SQLUser) Validate() error {
	return nil
}
