package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/go2hal/hal"
)

// UsersListHandler handles requests for users
func (h *Handler) UsersListHandler(w http.ResponseWriter, r *http.Request) {
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

	rows, err := h.db.Query(`SELECT Username
		FROM Users`)

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
		embeddedUser.Data()["username"] = username
		embeddedUsers = append(embeddedUsers, embeddedUser)
	}
	users, _ := hal.NewResourceRelation("users")
	users.SetResources(embeddedUsers)
	root.AddResource(users)

	w.Write(JSONify(root))
}

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
