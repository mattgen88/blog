package models

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"time"
)

// Article is an interface for describing articles
type Article interface {
	Exists() bool
	Populate() error
	Save() error
	Validate() error
}

// SQLArticle is a SQL backed Article
type SQLArticle struct {
	ID        int          `json:"id"`
	Author    *SQLUser     `json:"author"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	Date      *time.Time   `json:"date"`
	Slug      string       `json:"slug"`
	Category  *SQLCategory `json:"category"`
	Db        *sql.DB      `json:"-"`
	populated bool
	dirty     bool
	exists    bool
}

// NewSQLArticle returns a new instance of SQLArticle backed by a database
func NewSQLArticle(slug string, Db *sql.DB) *SQLArticle {
	p := &SQLArticle{
		Slug: slug,
		Db:   Db,
	}

	if p.Exists() {
		p.Populate()
	}

	return p
}

// ArticleListByCategory returns an article list by category, imagine that.
func ArticleListByCategory(categoryId int, Db *sql.DB) []*SQLArticle {
	var articles []*SQLArticle

	rows, err := Db.Query(`SELECT ArticleId, Title, Slug, Date, Users.Username
		FROM Articles
		JOIN Users on Users.UserId = Articles.Author
		WHERE Articles.Category = ?`, categoryId)

	if err != nil {
		log.Println(err)
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
			continue
		}

		article := &SQLArticle{
			Db:     Db,
			ID:     articleId,
			Title:  title,
			Slug:   slug,
			Date:   date,
			Author: NewSQLUser(author, Db),
		}

		articles = append(articles, article)

	}
	return articles
}

// ArticleList is a list of articles
func ArticleList(Db *sql.DB) []*SQLArticle {
	var articles []*SQLArticle

	rows, err := Db.Query(`SELECT ArticleId, Title, Slug, Date, Users.Username, Name, Body
		FROM Articles
		JOIN Category on Category.CategoryId = Articles.Category
		JOIN Users on Users.UserId = Articles.Author`)

	if err != nil {
		log.Println("Error querying for all articles", err)
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
			body      string
		)

		if err := rows.Scan(&articleId, &title, &slug, &date, &author, &category, &body); err != nil {
			continue
		}

		article := &SQLArticle{
			Db:       Db,
			ID:       articleId,
			Title:    title,
			Slug:     slug,
			Body:     body,
			Date:     date,
			Category: NewSQLCategory(category, Db),
			Author:   NewSQLUser(author, Db),
		}

		articles = append(articles, article)

	}
	return articles
}

// Exists determines whether or not the given post, by slug, exists
func (p *SQLArticle) Exists() bool {
	if p.exists {
		return true
	}
	var count int
	err := p.Db.QueryRow(`SELECT COUNT(*) FROM Articles WHERE Slug = ?`, p.Slug).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	p.exists = count > 0
	return count > 0
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

	err := p.Db.QueryRow(`SELECT ArticleId, Title, Users.Username, Body, Date, Slug, Name
	FROM Articles
	JOIN Category ON Articles.Category = Category.CategoryID
	JOIN Users ON Articles.Author = Users.UserId
	WHERE Slug = ?`, p.Slug).Scan(&p.ID, &p.Title, &author, &p.Body, &p.Date, &p.Slug, &category)

	if err != nil {
		log.Println("Item does not exist")
		return DNEError
	}

	p.Author = NewSQLUser(author, p.Db)
	p.Category = NewSQLCategory(category, p.Db)

	p.populated = true
	return nil
}

// Save the properties of the article into the database
func (p *SQLArticle) Save() error {
	var err error
	var query string

	err = p.Validate()
	if err != nil {
		// Validation error
		return err
	}
	if !p.Exists() {
		query = "INSERT INTO Articles (Title, Author, Body, Date, Slug, Category) VALUES (?, ?, ?, ?, ?, ?)"
		result, err := p.Db.Exec(query, p.Title, p.Author.ID, p.Body, p.Date, p.Slug, p.Category.ID)
		id, err := result.LastInsertId()
		if err != nil {
			return SaveError
		}
		p.ID = int(id)
	} else {
		query = "UPDATE Articles SET Title = ?, Author = ?, Body = ?, Date = ?, Slug = ?, Category = ? WHERE ArticleId = ?"
		_, err = p.Db.Exec(query, p.Title, p.Author.ID, p.Body, p.Date, p.Slug, p.Category.ID, p.ID)
	}

	if err != nil {
		return SaveError
	}

	return nil
}

var slugRegexp = regexp.MustCompile(`[[:alpha:]]+`)

// Validate the properties of model
func (p *SQLArticle) Validate() error {
	// Check each of the properties and validate them
	p.Category.Db = p.Db
	p.Category.Populate()
	p.Author.Db = p.Db
	p.Author.Populate()

	match := slugRegexp.MatchString(p.Slug)
	if !match {
		log.Println("Failed to pass regex", p.Slug)
		return ValidationError
	}

	return nil
}
