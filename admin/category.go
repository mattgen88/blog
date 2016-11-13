package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"

	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/blog/util"
)

// CreateCategoryHandler allows for creating new categories
func (a *Handler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)
	root.Data()["test"] = "testing"

	w.Write(util.JSONify(root))
}

// ReplaceCategoryHandler should take posts of categories and save them to the database
// after checking for possible problems
func (a *Handler) ReplaceCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Set up our hal resource
	root := hal.NewResourceObject()

	category := mux.Vars(r)["category"]

	// Fetch the requested article
	model := models.NewSQLCategory(category, a.db)
	model.Populate()

	// Unpack posted data into model
	err := json.NewDecoder(r.Body).Decode(&model)

	if err != nil {
		// parse error
		root.Data()["error"] = fmt.Sprintf("%s", ParseError)
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
	root.Data()["id"] = model.ID
	root.Data()["category"] = model.Name

	link := &hal.LinkObject{Href: "/categories/" + model.Name}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	w.Write(util.JSONify(root))
}
