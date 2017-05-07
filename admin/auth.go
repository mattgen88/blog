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

type userDataKey string
type role string
type userData struct {
	role     role
	username string
}

// Claims holds claims for a token
type Claims struct {
	jwt.StandardClaims
	claims map[string]string
}

// Auth handles request to authenticate and will issue a JWT
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	if r.Method != http.MethodPost {
		root.Data["error"] = "Please POST credentials."
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
		root.Data["error"] = "Missing required fields."
		root.Data["result"] = false

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
		root.Data["error"] = "Unable to authenticate. Check that credentials are correct."
		root.Data["result"] = false

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

	accessExpires := now.Add(time.Minute * 5)
	refreshExpires := now.Add(time.Hour * 24)

	// Create the Claims
	accessClaims := Claims{
		claims: map[string]string{
			"username": model.Username,
			"role":     model.Role,
		},
	}
	accessClaims.ExpiresAt = accessExpires.Unix()
	accessClaims.Issuer = "test"

	refreshClaims := Claims{
		claims: map[string]string{
			"username": model.Username,
		},
	}
	refreshClaims.ExpiresAt = refreshExpires.Unix()
	refreshClaims.Issuer = "test"

	accessCookie, accessErr := createJwt("access.jwt", accessExpires, &accessClaims, h.jwtKey)
	refreshCookie, refreshErr := createJwt("refresh.jwt", refreshExpires, &refreshClaims, h.jwtKey)
	if accessErr != nil {
		root.Data["err"] = fmt.Sprintf("%s", accessErr)
		root.Data["result"] = false
	} else if refreshErr != nil {
		root.Data["err"] = fmt.Sprintf("%s", refreshErr)
		root.Data["result"] = false
	} else {
		root.Data["result"] = true
		root.Data["access_expires"] = accessExpires.Unix()
		root.Data["refresh_expires"] = refreshExpires.Unix()
		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)
		root.AddLink("refresh", &haljson.Link{Href: "/refresh"})
	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

func createJwt(name string, expire time.Time, claims jwt.Claims, jwtKey string) (*http.Cookie, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		Expires:  expire,
	}

	return &cookie, nil
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
