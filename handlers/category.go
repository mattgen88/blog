package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/haljson"
)

// CategoryHandler handles requests for categories
func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	c := mux.Vars(r)["category"]

	category := models.NewSQLCategory(c, h.db)

	root.Data["id"] = category.ID

	categories := models.ArticleListByCategory(category.ID, h.db)

	for _, article := range categories {

		href := fmt.Sprintf("/articles/%s", article.Slug)
		embeddedArticle := haljson.NewResource()
		embeddedArticle.Self(href)
		embeddedArticle.Data["title"] = article.Title
		embeddedArticle.Data["author"] = article.Author.Username
		embeddedArticle.Data["date"] = article.Date
		root.AddEmbed("articles", embeddedArticle)
	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// CategoryListHandler requests a list of categories
func (h *Handler) CategoryListHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	var categories []string

	for _, category := range models.CategoryList(h.db) {

		href := fmt.Sprintf("/categories/%s", category.Name)

		embeddedCategory := haljson.NewResource()

		embeddedCategory.Self(href)

		embeddedCategory.Data["name"] = category.Name
		root.AddEmbed("categories", embeddedCategory)
		categories = append(categories, category.Name)
	}
	root.Data["categories"] = categories

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
