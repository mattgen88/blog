package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattgen88/blog/models"
	"github.com/pmoule/go2hal/hal"
)

// CategoryHandler handles requests for categories
// @TODO: return list of articles in that category and embed them
func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	c := mux.Vars(r)["category"]

	category := models.NewSQLCategory(c, h.db)

	root.Data()["id"] = category.GetID()

	w.Write(JSONify(root))
}

// ArticleListHandler handles requests for articles
func (h *Handler) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	// @TODO: Move this off to the model...
	rows, err := h.db.Query(`SELECT ArticleId, Title, Body, Date, Slug, Category.Name, Users.Username
		FROM Articles
		JOIN Category on Category.CategoryId = Articles.Category
		JOIN Users on Users.UserId = Articles.Author`)

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var embeddedArticles []hal.Resource

	for rows.Next() {
		var (
			postId   int
			title    string
			body     string
			date     string
			slug     string
			category string
			author   string
		)

		if err := rows.Scan(&postId, &title, &body, &date, &slug, &category, &author); err != nil {
			fmt.Println(err)
			continue
		}

		href := "/articles/" + category + "/" + slug + "/"
		selfLink, err := hal.NewLinkObject(href)

		if err != nil {
			fmt.Println(err)
		}

		self = hal.NewSelfLinkRelation()
		self.SetLink(selfLink)

		embeddedArticle := hal.NewResourceObject()
		embeddedArticle.AddLink(self)
		embeddedArticle.Data()["title"] = title
		embeddedArticle.Data()["author"] = author
		embeddedArticle.Data()["category"] = category
		embeddedArticle.Data()["date"] = date
		embeddedArticles = append(embeddedArticles, embeddedArticle)
	}

	articles, _ := hal.NewResourceRelation("articles")
	articles.SetResources(embeddedArticles)
	root.AddResource(articles)

	w.Write(JSONify(root))
}

// ArticleHandler handles requests for articles
func (h *Handler) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["category"] = mux.Vars(r)["category"]
	root.Data()["id"] = mux.Vars(r)["id"]

	w.Write(JSONify(root))
}
