package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattgen88/blog/models"
	"github.com/pmoule/go2hal/hal"
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

		href := fmt.Sprintf("/articles/%s/%s", c, article.Slug)
		selfLink, err := hal.NewLinkObject(href)

		if err != nil {
			fmt.Println(err)
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

	w.Write(JSONify(root))
}

// ArticleListHandler handles requests for articles
func (h *Handler) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	var embeddedArticles []hal.Resource

	for _, article := range models.ArticleList(h.db) {

		href := fmt.Sprintf("/articles/%s/%s", article.Category.Name, article.Slug)
		selfLink, err := hal.NewLinkObject(href)

		if err != nil {
			fmt.Println(err)
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

	w.Write(JSONify(root))
}

// ArticleHandler handles requests for articles
func (h *Handler) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(mux.Vars(r)["id"])
	article := models.NewSQLArticle(mux.Vars(r)["id"], h.db)

	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["category"] = mux.Vars(r)["category"]
	root.Data()["id"] = mux.Vars(r)["id"]
	root.Data()["body"] = article.GetBody()

	w.Write(JSONify(root))
}
