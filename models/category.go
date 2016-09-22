package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Category in an interface for categories
type Category interface {
	GetName() string
	GetId() int
	Populate() error
}

// SQLCategory is a Category backed by SQL
type SQLCategory struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name"`
	Db        *sql.DB `json:"-"`
	populated bool
	dirty     bool
}

// NewSQLCategory creates a SQLCategory instance configured with a connection
func NewSQLCategory(name string, db *sql.DB) *SQLCategory {
	c := &SQLCategory{
		Db:   db,
		Name: name,
	}

	if c.Exists() {
		err := c.Populate()
		if err != nil {
			log.Println(err)
		}
	}

	return c
}

// GetID returns the ID of a category
func (c *SQLCategory) GetID() int {
	return c.ID
}

// Exists check if the category exists
func (c *SQLCategory) Exists() bool {
	var count int
	err := c.Db.QueryRow(`SELECT COUNT(*)
	FROM Category
	WHERE Name = ?`, c.Name).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}

	return true
}

// Populate the model with data from the database
func (c *SQLCategory) Populate() error {
	if !c.Exists() {
		return errors.New("Instance does not exist")
	}

	if c.populated {
		return errors.New("Model already populated")
	}

	if c.dirty {
		return errors.New("Model dirty")
	}

	// Fetch data and populate
	err := c.Db.QueryRow(`SELECT CategoryId
	FROM Category
	WHERE Name = ?`, c.Name).Scan(&c.ID)

	if err != nil {
		return errors.New("Unknown error occurred: " + fmt.Sprintf("%s", err))
	}

	c.populated = true
	return nil
}
