package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"
)

// CategoryHandler handles requests for categories
func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["request"] = mux.Vars(r)["category"]

	w.Write(JSONify(root))
}

// ArticleListHandler handles requests for articles
func (h *Handler) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	rows, err := h.db.Query(`SELECT PostId, Title, Body, Date, Slug, Category.Name, Users.Username
		FROM Posts
		JOIN Category on Category.CategoryId = Posts.Category
		JOIN Users on Users.UserId = Posts.Author`)
	
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var embeddedPosts []hal.Resource

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

		embeddedPost := hal.NewResourceObject()
		embeddedPost.AddLink(self)
		embeddedPost.Data()["title"] = title
		embeddedPost.Data()["author"] = author
		embeddedPost.Data()["category"] = category
		embeddedPost.Data()["date"] = date
		embeddedPosts = append(embeddedPosts, embeddedPost)
	}
	posts, _ := hal.NewResourceRelation("posts")
	posts.SetResources(embeddedPosts)
	root.AddResource(posts)

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
