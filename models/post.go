package models

import (
	"database/sql"
	"errors"
	"time"
)

// Post is an interface for describing posts
type Post interface {
	Exists() bool
	Populate() error
	GetSlug() string
	GetAuthor() *User
	GetTitle() string
	GetBody() string
	GetDate() *time.Time
	GetCategory() *SQLCategory
}

// SQLPost is a SQL backed Post
type SQLPost struct {
	db        *sql.DB
	ID        int
	Author    *SQLUser
	Title     string
	Body      string
	Date      *time.Time
	Slug      string
	Category  *SQLCategory
	populated bool
	dirty     bool
}

// Return a new instance of SQLPost backed by a database
func NewSQLPost(slug string, db *sql.DB) *SQLPost {
	p := &SQLPost{
		db: db,
	}
	if p.Exists() {
		p.Populate()
	}
	return p
}

// GetAuthor returns the User who authored the post
func (p *SQLPost) GetAuthor() *SQLUser {
	if p.populated {
		return p.Author
	}
	return nil
}

// GetTitle returns the title of the post
func (p *SQLPost) GetTitle() string {
	if p.populated {
		return p.Title
	}
	return ""
}

// GetBody returns the body of the post
func (p *SQLPost) GetBody() string {
	if p.populated {
		return p.Body
	}
	return ""
}

// GetDate returns the date the post was authored
func (p *SQLPost) GetDate() *time.Time {
	if p.populated {
		return p.Date
	}
	return nil
}

// GetSlug returns the post's slug to uniquely identify it in URLs
func (p *SQLPost) GetSlug() string {
	if p.populated {
		return p.Slug
	}
	return ""
}

// GetCategory returns the Category for the post
func (p *SQLPost) GetCategory() *SQLCategory {
	if p.populated {
		return p.Category
	}
	return nil
}

// Exists determines whether or not the given post, by slug, exists
func (p *SQLPost) Exists() bool {
	var count int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM Posts WHERE Slug = ?`, p.Slug).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}

// Populate populates the model with data from the database
func (p *SQLPost) Populate() error {
	if !p.Exists() {
		return errors.New("instance does not exist")
	}
	if p.populated {
		return errors.New("Model already populated")
	}
	if p.dirty {
		return errors.New("Model dirty")
	}

	var (
		author   string
		category string
	)

	err := p.db.QueryRow(`SELECT PostId, Title, Users.Username, Body, Date, Slug, Category.Name
	FROM Posts, Category
	JOIN Category ON Posts.Category = Category.CategoryID
	JOIN Users ON Posts.Author = Users.UserId
	WHERE Slug = ?`).Scan(&p.ID, &p.Title, &author, &p.Body, &p.Date, &p.Slug, &category)
	if err != nil {
		return errors.New("Unknown error occurred")
	}

	p.Author = NewSQLUser(author, p.db)
	p.Category = NewSQLCategory(category, p.db)

	p.populated = true
	return nil
}
