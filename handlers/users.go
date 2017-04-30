package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattgen88/haljson"
)

// UsersListHandler handles requests for users
func (h *Handler) UsersListHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self(r.URL.Path)

	rows, err := h.db.Query(`SELECT Username
		FROM Users`)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var username string
		if scanErr := rows.Scan(&username); scanErr != nil {
			log.Println(scanErr)
			continue
		}

		href := "/users/" + username

		embeddedUser := haljson.NewResource()
		embeddedUser.Self(href)
		embeddedUser.Data["username"] = username
		root.AddEmbed("users", embeddedUser)

	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// UserHandler handles requests for users
func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	root.Data["username"] = mux.Vars(r)["id"]

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
