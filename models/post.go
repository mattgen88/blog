package models

import (
	"database/sql"
	"errors"
	"time"
)

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

func NewSQLPost(slug string, db *sql.DB) *SQLPost {
	p := &SQLPost{
		db: db,
	}
	if p.Exists() {
		p.Populate()
	}
	return p
}

func (p *SQLPost) GetAuthor() *SQLUser {
	if p.populated {
		return p.Author
	}
	return nil
}
func (p *SQLPost) GetTitle() string {
	if p.populated {
		return p.Title
	}
	return ""
}
func (p *SQLPost) GetBody() string {
	if p.populated {
		return p.Body
	}
	return ""
}
func (p *SQLPost) GetDate() *time.Time {
	if p.populated {
		return p.Date
	}
	return nil
}
func (p *SQLPost) GetSlug() string {
	if p.populated {
		return p.Slug
	}
	return ""
}
func (p *SQLPost) GetCategory() *SQLCategory {
	if p.populated {
		return p.Category
	}
	return nil
}

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
