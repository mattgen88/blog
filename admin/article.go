package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/haljson"
)

// CreateArticleHandler allows for creating new articles
func (a *Handler) CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	// @TODO: Fix author and category lookup
	// Set up our hal resource
	root := haljson.NewResource()

	root.Self(r.URL.Path)

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
		root.Data["error"] = fmt.Sprintf("%s", ErrParse)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
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
		root.Data["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	// Write the model out
	root.Data["title"] = model.Title
	root.Data["author"] = model.Author
	root.Data["body"] = model.Body
	root.Data["slug"] = model.Slug
	root.Data["date"] = model.Date
	root.Data["id"] = model.ID
	root.Data["category"] = model.Category

	json, marshalErr := json.Marshal(root)
	if marshalErr != nil {
		log.Println(marshalErr)
		return
	}
	w.Write(json)
}

// ReplaceArticleHandler should take posts of articles and save them to the database
// after checking for possible problems
func (a *Handler) ReplaceArticleHandler(w http.ResponseWriter, r *http.Request) {

	// Set up our hal resource
	root := haljson.NewResource()

	root.Self(r.URL.Path)

	// Get the slug of the post we're dealing with
	slug := mux.Vars(r)["id"]

	// Fetch the requested article
	model := models.NewSQLArticle(slug, a.db)
	err := model.Populate()
	if err != nil {
		root.Data["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	// Unpack posted data into model
	err = json.NewDecoder(r.Body).Decode(&model)

	if err != nil {
		// parse error
		root.Data["error"] = fmt.Sprintf("%s", ErrParse)
		w.WriteHeader(http.StatusBadRequest)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	err = model.Save()

	if err != nil {
		root.Data["error"] = fmt.Sprintf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	// Write the model out
	root.Data["title"] = model.Title
	root.Data["author"] = model.Author
	root.Data["body"] = model.Body
	root.Data["slug"] = model.Slug
	root.Data["date"] = model.Date
	root.Data["id"] = model.ID
	root.Data["category"] = model.Category

	json, marshalErr := json.Marshal(root)
	if marshalErr != nil {
		log.Println(marshalErr)
		return
	}
	w.Write(json)
}
