package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	log.Println("auth request received", r.FormValue("username"))
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	// @TODO look at r.Method for GET vs POST and provide some information about correctly authenticating

	model := models.NewSQLUser(r.FormValue("username"), h.db)
	err := model.Populate()
	if err != nil || !model.Authenticate(r.FormValue("password")) {
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
	} else {
		root.Data["jwt"] = tokenString
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Secure:   true,
			HttpOnly: true,
			Expires:  expires,
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
	root.Data["test"] = "create user"

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)

}
