package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// CreateUserHandler allows for the creation of users
func (a *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["test"] = "create user"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// ReplaceUserHandler allows for the modification of users
func (a *Handler) ReplaceUserHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["test"] = "replace user"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
