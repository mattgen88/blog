package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Category in an interface for categories
type Category interface {
	Populate() error
	Save() error
	Validate() error
}

// SQLCategory is a Category backed by SQL
type SQLCategory struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name"`
	Db        *sql.DB `json:"-"`
	populated bool
	dirty     bool
	exists    bool
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

// CategoryList is a list of categories
func CategoryList(Db *sql.DB) []*SQLCategory {
	var categories []*SQLCategory

	rows, err := Db.Query(`SELECT CategoryId, Name from Category`)

	if err != nil {
		log.Println("Error querying for all categories", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			categoryId int
			name       string
		)

		if err := rows.Scan(&categoryId, &name); err != nil {
			log.Println(err)
			continue
		}

		category := &SQLCategory{
			Db:   Db,
			ID:   categoryId,
			Name: name,
		}

		categories = append(categories, category)

	}
	return categories
}

// Exists check if the category exists
func (c *SQLCategory) Exists() bool {
	if c.exists {
		return true
	}
	var count int
	err := c.Db.QueryRow(`SELECT COUNT(*)
	FROM Category
	WHERE Name = ?`, c.Name).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	c.exists = count > 0
	return count > 0
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

// Save the properties of the category into the database
func (c *SQLCategory) Save() error {
	var err error
	var query string

	err = c.Validate()
	if err != nil {
		// Validation error
		return err
	}
	if !c.Exists() {
		log.Println("Creating new category")
		query = `INSERT INTO Category ("Name") VALUES (?)`
		_, err = c.Db.Exec(query, c.Name)
	} else {
		log.Println("Overwriting existing category")
		query = `UPDATE Category SET "Name" = ? WHERE "CategoryId" = ?`
		_, err = c.Db.Exec(query, c.Name, c.ID)
	}

	if err != nil {
		return SaveError
	}

	return nil
}

// Validate the properties of category
func (c *SQLCategory) Validate() error {
	log.Println(c.Name)
	match, err := regexp.MatchString(`[a-zA-Z0-9\-_]+`, strings.TrimSpace(c.Name))
	if err != nil {
		log.Println(err)
		return ValidationError
	}
	if !match {
		log.Println("No match")
		return ValidationError
	}
	return nil
}
