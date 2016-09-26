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

		href := fmt.Sprintf("/articles/%s", c, article.Slug)
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
