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

// ArticleListHandler handles requests for articles
func (h *Handler) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	var embeddedArticles []hal.Resource

	for _, article := range models.ArticleList(h.db) {

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
		embeddedArticle.Data()["category"] = article.Category.Name
		embeddedArticle.Data()["slug"] = article.Slug

		var trunc int
		if len(article.Body) <= 100 {
			trunc = len(article.Body)
		} else {
			trunc = 100
		}
		embeddedArticle.Data()["description"] = article.Body[0:trunc]

		embeddedArticles = append(embeddedArticles, embeddedArticle)
	}

	articles, _ := hal.NewResourceRelation("articles")
	articles.SetResources(embeddedArticles)

	root.AddResource(articles)

	w.Write(util.JSONify(root))
}

// ArticleHandler handles requests for articles
func (h *Handler) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(mux.Vars(r)["id"])
	article := models.NewSQLArticle(mux.Vars(r)["id"], h.db)

	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.Data()["article"] = article

	root.AddLink(self)
	root.Data()["body"] = article.Body
	root.Data()["title"] = article.Title
	root.Data()["author"] = article.Author.Username
	root.Data()["date"] = article.Date
	root.Data()["category"] = article.Category.Name
	root.Data()["slug"] = article.Slug

	w.Write(util.JSONify(root))
}
