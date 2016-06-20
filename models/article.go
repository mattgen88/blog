package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Article is an interface for describing articles
type Article interface {
	Exists() bool
	Populate() error
	GetSlug() string
	GetAuthor() *User
	GetTitle() string
	GetBody() string
	GetDate() *time.Time
	GetCategory() *SQLCategory
}

// SQLArticle is a SQL backed Article
type SQLArticle struct {
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

// Return a new instance of SQLArticle backed by a database
func NewSQLArticle(slug string, db *sql.DB) *SQLArticle {
	p := &SQLArticle{
		Slug: slug,
		db: db,
	}

	if p.Exists() {
		p.Populate()
	}

	return p
}

func ArticleListByCategory(categoryId int, db *sql.DB) []*SQLArticle {
	var articles []*SQLArticle

	rows, err := db.Query(`SELECT ArticleId, Title, Slug, Date, Users.Username
		FROM Articles
		JOIN Users on Users.UserId = Articles.Author
		WHERE Articles.Category = ?`, categoryId)

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			articleId int
			title     string
			date      *time.Time
			slug      string
			author    string
		)

		if err := rows.Scan(&articleId, &title, &slug, &date, &author); err != nil {
			fmt.Println(err)
			continue
		}

		article := &SQLArticle{
			db:       db,
			ID:       articleId,
			Title:    title,
			Slug:     slug,
			Date:     date,
			Author:   NewSQLUser(author, db),
		}

		articles = append(articles, article)

	}
	return articles
}

func ArticleList(db *sql.DB) []*SQLArticle {
	var articles []*SQLArticle

	// @TODO: Move this off to the model...
	rows, err := db.Query(`SELECT ArticleId, Title, Slug, Date, Users.Username, Name
		FROM Articles
		JOIN Category on Category.CategoryId = Articles.Category
		JOIN Users on Users.UserId = Articles.Author`)

	if err != nil {
		fmt.Println("Error querying for all articles", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			articleId int
			title     string
			date      *time.Time
			slug      string
			author    string
			category  string
		)

		if err := rows.Scan(&articleId, &title, &slug, &date, &author, &category); err != nil {
			fmt.Println(err)
			continue
		}

		article := &SQLArticle{
			db:     db,
			ID:     articleId,
			Title:  title,
			Slug:   slug,
			Date:   date,
			Category: NewSQLCategory(category, db),
			Author: NewSQLUser(author, db),
		}

		articles = append(articles, article)

	}
	return articles
}

// GetAuthor returns the User who authored the post
func (p *SQLArticle) GetAuthor() *SQLUser {
	if p.populated {
		return p.Author
	}

	return nil
}

// GetTitle returns the title of the post
func (p *SQLArticle) GetTitle() string {
	if p.populated {
		return p.Title
	}

	return ""
}

// GetBody returns the body of the post
func (p *SQLArticle) GetBody() string {
	if p.populated {
		return p.Body
	}

	return ""
}

// GetDate returns the date the post was authored
func (p *SQLArticle) GetDate() *time.Time {
	if p.populated {
		return p.Date
	}

	return nil
}

// GetSlug returns the post's slug to uniquely identify it in URLs
func (p *SQLArticle) GetSlug() string {
	if p.populated {
		return p.Slug
	}

	return ""
}

// GetCategory returns the Category for the post
func (p *SQLArticle) GetCategory() *SQLCategory {
	if p.populated {
		return p.Category
	}

	return nil
}

// Exists determines whether or not the given post, by slug, exists
func (p *SQLArticle) Exists() bool {
	var count int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM Articles WHERE Slug = ?`, p.Slug).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows for slug " + p.Slug)
			return false
		}
		fmt.Println("Some other error", err)
		return false
	}
	fmt.Println("Found " + p.Slug);
	return true
}

// Populate populates the model with data from the database
func (p *SQLArticle) Populate() error {
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

	err := p.db.QueryRow(`SELECT ArticleId, Title, Users.Username, Body, Date, Slug, Name
	FROM Articles
	JOIN Category ON Articles.Category = Category.CategoryID
	JOIN Users ON Articles.Author = Users.UserId
	WHERE Slug = ?`, p.Slug).Scan(&p.ID, &p.Title, &author, &p.Body, &p.Date, &p.Slug, &category)

	if err != nil {
		fmt.Println("Unknown error", err)
		return errors.New("Unknown error occurred")
	}

	fmt.Println(p)

	p.Author = NewSQLUser(author, p.db)
	p.Category = NewSQLCategory(category, p.db)

	p.populated = true
	return nil
}