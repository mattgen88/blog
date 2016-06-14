package models

import (
	"database/sql"
	"errors"
)

// Category in an interface for categories
type Category interface {
	GetName() string
	GetId() int
}

// SQLCategory is a Category backed by SQL
type SQLCategory struct {
	db        *sql.DB
	ID        int
	Name      string
	populated bool
	dirty     bool
}

// NewSQLCategory creates a SQLCategory instance configured with a connection
func NewSQLCategory(name string, db *sql.DB) *SQLCategory {
	c := &SQLCategory{
		db:   db,
		Name: name,
	}
	if c.Exists() {
		c.Populate()
	}
	return c
}

// Exists check if the category exists
func (c *SQLCategory) Exists() bool {
	var count int
	err := c.db.QueryRow(`SELECT COUNT(*) FROM Category WHERE Name = ?`, c.Name).Scan(&count)

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
	err := c.db.QueryRow(`SELECT CategoryId, Name
	FROM Categories
	WHERE Name = ?`, c.Name).Scan(c.ID)

	if err != nil {
		return errors.New("Unknown error occurred")
	}
	c.populated = true
	return nil
}