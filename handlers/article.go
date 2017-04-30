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

// ArticleListHandler handles requests for articles
func (h *Handler) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	for _, article := range models.ArticleList(h.db) {

		href := fmt.Sprintf("/articles/%s", article.Slug)

		embeddedArticle := haljson.NewResource()
		embeddedArticle.Self(href)
		embeddedArticle.Data["title"] = article.Title
		embeddedArticle.Data["author"] = article.Author.Username
		embeddedArticle.Data["date"] = article.Date
		embeddedArticle.Data["category"] = article.Category.Name
		embeddedArticle.Data["slug"] = article.Slug

		var trunc int
		if len(article.Body) <= 100 {
			trunc = len(article.Body)
		} else {
			trunc = 100
		}
		embeddedArticle.Data["description"] = article.Body[0:trunc]
		root.AddEmbed("articles", embeddedArticle)
	}
	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// ArticleHandler handles requests for articles
func (h *Handler) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	article := models.NewSQLArticle(mux.Vars(r)["id"], h.db)

	root.Data["article"] = article
	root.Data["body"] = article.Body
	root.Data["title"] = article.Title
	root.Data["author"] = article.Author.Username
	root.Data["date"] = article.Date
	root.Data["category"] = article.Category.Name
	root.Data["slug"] = article.Slug

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
