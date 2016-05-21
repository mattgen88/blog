package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"
)

// Handler provides various http handlers
type Handler struct { // {{{
	r  *mux.Router
	db *sql.DB
} // }}}

// New returns a configured handler struct
func New(r *mux.Router, db *sql.DB) *Handler { // {{{
	return &Handler{r, db}
} // }}}

// JSONify the resource
func JSONify(root hal.Resource) []byte { // {{{

	encoder := new(hal.Encoder)
	bytes, err := encoder.ToJSON(root)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return bytes
} // }}}

// RootHandler handles requests for the root of the API
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) { // {{{
	root := hal.NewResourceObject()

	self := hal.NewSelfLinkRelation()
	self.SetLink(&hal.LinkObject{Href: r.URL.Path})

	root.AddLink(self)

	// Users
	users, err := hal.NewLinkRelation("Users")
	if err != nil {
		// Error creating link relation
		fmt.Println(err)
		return
	}
	users.SetLink(&hal.LinkObject{Href: "/users/"})

	root.AddLink(users)

	// Article
	article, err := hal.NewLinkRelation("Article")
	if err != nil {
		// Error creating link relation
		fmt.Println(err)
		return
	}
	article.SetLink(&hal.LinkObject{Href: "/articles/{category}/{id:[a-zA-Z-_]+}"})

	root.AddLink(article)

	// Article Category
	category, err := hal.NewLinkRelation("Article Category")
	if err != nil {
		// Error creating link relation
		fmt.Println(err)
		return
	}
	category.SetLink(&hal.LinkObject{Href: "/articles/{category}"})

	root.AddLink(category)

	// Write it out
	w.Write(JSONify(root))
} // }}}

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
func (h *Handler) ArticleListHandler(w http.ResponseWriter, r *http.Request) { // {{{
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	rows, err := h.db.Query("SELECT PostId, Title, Body, Date, Slug, Category.Name, Users.Username from Posts JOIN Category on Category.CategoryId = Posts.Category JOIN Users on Users.UserId = Posts.Author")
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
} // }}}

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

// UsersListHandler handles requests for users
func (h *Handler) UsersListHandler(w http.ResponseWriter, r *http.Request) { // {{{
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	rows, err := h.db.Query("SELECT Username from Users")
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var embeddedUsers []hal.Resource

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			fmt.Println(err)
			continue
		}

		href := "/users/" + username
		selfLink, err := hal.NewLinkObject(href)
		if err != nil {
			fmt.Println(err)
		}

		self = hal.NewSelfLinkRelation()
		self.SetLink(selfLink)

		embeddedUser := hal.NewResourceObject()
		embeddedUser.AddLink(self)
		embeddedUser.Data()["name"] = username
		embeddedUsers = append(embeddedUsers, embeddedUser)
	}
	users, _ := hal.NewResourceRelation("users")
	users.SetResources(embeddedUsers)
	root.AddResource(users)

	w.Write(JSONify(root))
} // }}}

// UserHandler handles requests for users
func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["username"] = mux.Vars(r)["id"]

	w.Write(JSONify(root))
}

// ErrorHandler handles requests for users
func (h *Handler) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["message"] = "Resource not found"

	w.WriteHeader(http.StatusNotFound)

	w.Write(JSONify(root))
}
