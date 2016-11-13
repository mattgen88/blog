package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/blog/util"
)

// CreateArticleHandler allows for creating new articles
func (a *Handler) CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	// @TODO: Fix author and category lookup
	// Set up our hal resource
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	// Get the slug of the post we're dealing with
	slug := mux.Vars(r)["id"]

	// Fetch the requested article
	model := models.NewSQLArticle(slug, a.db)

	// Unpack posted data into model
	err := json.NewDecoder(r.Body).Decode(&model)

	if err != nil {
		// parse error
		log.Println("Error parsing")
		log.Println(err)
		root.Data()["error"] = fmt.Sprintf("%s", ParseError)
		w.Write(util.JSONify(root))
		return
	}

	if model.Exists() {
		log.Println("Conflict")
		w.WriteHeader(http.StatusConflict)
		return
	}
	now := time.Now()
	model.Date = &now
	err = model.Save()

	if err != nil {
		log.Println("Error saving")
		root.Data()["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(util.JSONify(root))
		return
	}

	// Write the model out
	root.Data()["title"] = model.Title
	root.Data()["author"] = model.Author
	root.Data()["body"] = model.Body
	root.Data()["slug"] = model.Slug
	root.Data()["date"] = model.Date
	root.Data()["id"] = model.ID
	root.Data()["category"] = model.Category

	w.Write(util.JSONify(root))
}

// ReplaceArticleHandler should take posts of articles and save them to the database
// after checking for possible problems
func (a *Handler) ReplaceArticleHandler(w http.ResponseWriter, r *http.Request) {

	// Set up our hal resource
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	// Get the slug of the post we're dealing with
	slug := mux.Vars(r)["id"]

	// Fetch the requested article
	model := models.NewSQLArticle(slug, a.db)
	err := model.Populate()
	if err != nil {
		root.Data()["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(util.JSONify(root))
		return
	}

	// Unpack posted data into model
	err = json.NewDecoder(r.Body).Decode(&model)

	if err != nil {
		// parse error
		root.Data()["error"] = fmt.Sprintf("%s", ParseError)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(util.JSONify(root))
		return
	}

	err = model.Save()

	if err != nil {
		root.Data()["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(util.JSONify(root))
		return
	}

	// Write the model out
	root.Data()["title"] = model.Title
	root.Data()["author"] = model.Author
	root.Data()["body"] = model.Body
	root.Data()["slug"] = model.Slug
	root.Data()["date"] = model.Date
	root.Data()["id"] = model.ID
	root.Data()["category"] = model.Category

	w.Write(util.JSONify(root))
}
