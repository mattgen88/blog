package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/haljson"
)

// AuthClaims jwt claims
type AuthClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Auth handles request to authenticate and will issue a JWT
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	if r.Method != http.MethodPost {
		root.Data["error"] = "Please POST credentials"
		root.Data["required_fields"] = []string{"username", "password"}

		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(json)
		return
	}

	if r.FormValue("username") == "" || r.FormValue("password") == "" {
		root.Data["required_fields"] = []string{"username", "password"}

		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
		return
	}

	model := models.NewSQLUser(r.FormValue("username"), h.db)
	err := model.Populate()
	if err != nil || !model.Authenticate(r.FormValue("password")) {
		root.Data["error"] = "Unable to authenticate. Check that credentials are correct"
		w.WriteHeader(http.StatusForbidden)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	now := time.Now()
	expires := now.Add(time.Minute * 5)

	// Create the Claims
	claims := AuthClaims{
		model.Username,
		model.Role,
		jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			Issuer:    "test",
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(h.jwtKey))
	if err != nil {
		root.Data["err"] = fmt.Sprintf("%s", err)
		root.Data["result"] = false
	} else {
		root.Data["result"] = true
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Secure:   true,
			HttpOnly: true,
			Expires:  expires,
			Domain:   strings.Split(r.Host, ":")[0],
		}
		http.SetCookie(w, &cookie)
	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// AuthTest tests auth
func (h *Handler) AuthTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered auth test")
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["test"] = "Passed auth and loaded this"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)

}
