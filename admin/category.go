package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/haljson"
)

// CreateCategoryHandler allows for creating new categories
func (a *Handler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["test"] = "create category"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// ReplaceCategoryHandler should take posts of categories and save them to the database
// after checking for possible problems
func (a *Handler) ReplaceCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Set up our hal resource
	root := haljson.NewResource()

	category := mux.Vars(r)["category"]

	// Fetch the requested article
	model := models.NewSQLCategory(category, a.db)
	model.Populate()

	// Unpack posted data into model
	err := json.NewDecoder(r.Body).Decode(&model)

	if err != nil {
		// parse error
		root.Data["error"] = fmt.Sprintf("%s", ErrParse)
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
	root.Data["id"] = model.ID
	root.Data["category"] = model.Name

	root.Self("/categories/" + model.Name)

	json, marshalErr := json.Marshal(root)
	if marshalErr != nil {
		log.Println(marshalErr)
		return
	}
	w.Write(json)
}
