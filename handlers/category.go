package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/blog/util"
)

// CategoryHandler handles requests for categories
func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	c := mux.Vars(r)["category"]

	category := models.NewSQLCategory(c, h.db)

	root.Data()["id"] = category.GetID()

	categories := models.ArticleListByCategory(category.GetID(), h.db)

	var embeddedArticles []hal.Resource

	for _, article := range categories {

		href := fmt.Sprintf("/articles/%s", article.Slug)
		selfLink, err := hal.NewLinkObject(href)

		if err != nil {
			log.Println(err)
		}

		self = hal.NewSelfLinkRelation()
		self.SetLink(selfLink)

		embeddedArticle := hal.NewResourceObject()
		embeddedArticle.AddLink(self)
		embeddedArticle.Data()["title"] = article.Title
		embeddedArticle.Data()["author"] = article.Author.Username
		embeddedArticle.Data()["date"] = article.Date
		embeddedArticles = append(embeddedArticles, embeddedArticle)
	}

	articles, _ := hal.NewResourceRelation("articles")
	articles.SetResources(embeddedArticles)

	root.AddResource(articles)

	w.Write(util.JSONify(root))
}

// CategoryListHandler requests a list of categories
func (h *Handler) CategoryListHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	var embeddedCategories []hal.Resource

	for _, category := range models.CategoryList(h.db) {

		href := fmt.Sprintf("/categories/%s", category.Name)
		selfLink, err := hal.NewLinkObject(href)

		if err != nil {
			log.Println(err)
		}

		self = hal.NewSelfLinkRelation()
		self.SetLink(selfLink)

		embeddedCategory := hal.NewResourceObject()
		embeddedCategory.AddLink(self)

		embeddedCategory.Data()["name"] = category.Name

		embeddedCategories = append(embeddedCategories, embeddedCategory)
	}

	categories, _ := hal.NewResourceRelation("categories")
	categories.SetResources(embeddedCategories)

	root.AddResource(categories)

	w.Write(util.JSONify(root))
}
