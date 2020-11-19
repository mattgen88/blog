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
	Delete() error
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
func ArticleListByCategory(categoryID int, Db *sql.DB) []*SQLArticle {
	var articles []*SQLArticle

	rows, err := Db.Query(`SELECT "articleid", "title", "slug", "date", "users"."username"
		FROM "articles"
		JOIN "users" on "users"."userid" = "articles"."author"
		WHERE "articles"."category" = $1
		ORDER BY "date" desc`, categoryID)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			articleID int
			title     string
			date      *time.Time
			slug      string
			author    string
		)

		if err := rows.Scan(&articleID, &title, &slug, &date, &author); err != nil {
			continue
		}

		article := &SQLArticle{
			Db:     Db,
			ID:     articleID,
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

	rows, err := Db.Query(`SELECT "articleid", "title", "slug", "date", "users"."username", "name", "body"
		FROM "articles"
		JOIN "category" on "category"."categoryId" = "articles"."category"
		JOIN "users" on "users"."userid" = "articles"."author"
		ORDER BY "date" DESC`)

	if err != nil {
		log.Println("Error querying for all articles", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			articleID int
			title     string
			date      *time.Time
			slug      string
			author    string
			category  string
			body      string
		)

		if err := rows.Scan(&articleID, &title, &slug, &date, &author, &category, &body); err != nil {
			continue
		}

		article := &SQLArticle{
			Db:       Db,
			ID:       articleID,
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
	err := p.Db.QueryRow(`SELECT COUNT(*) FROM "articles" WHERE "slug" = $1`, p.Slug).Scan(&count)
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

	err := p.Db.QueryRow(`SELECT "articleid", "title", "users"."username", "body", "date", "slug", "name"
	FROM "articles"
	JOIN "category" ON "articles"."category" = "category"."categoryid"
	JOIN "users" ON "articles"."author" = "users"."userId"
	WHERE "slug" = $1`, p.Slug).Scan(&p.ID, &p.Title, &author, &p.Body, &p.Date, &p.Slug, &category)

	if err != nil {
		log.Println("Item does not exist", ErrDoesNotExist)
		return ErrDoesNotExist
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
		query = `INSERT INTO "articles" ("title", "author", "body", "date", "slug", "category") VALUES ($1, $2, $3, $4, $5, $6)`
		result, err := p.Db.Exec(query, p.Title, p.Author.ID, p.Body, p.Date, p.Slug, p.Category.ID)
		id, err := result.LastInsertId()
		if err != nil {
			return ErrSave
		}
		p.ID = int(id)
	} else {
		query = `UPDATE "articles" SET "title" = $1, "author" = $2, "body" = $3, "date" = $4, "slug" = $5, "category" = $6 WHERE "articleid" = $7`
		_, err = p.Db.Exec(query, p.Title, p.Author.ID, p.Body, p.Date, p.Slug, p.Category.ID, p.ID)
	}

	if err != nil {
		log.Println("Failed to save article", ErrSave)
		return ErrSave
	}

	return nil
}

// Delete the requested article
func (p *SQLArticle) Delete() error {
	var err error
	var query string

	if !p.Exists() {
		return ErrDoesNotExist
	}
	query = `DELETE FROM "articles" WHERE "slug" = $1`
	_, err = p.Db.Exec(query, p.Slug)

	if err != nil {
		log.Println("Failed to delete article", err)
		return ErrDelete
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
		return ErrValidation
	}

	return nil
}
