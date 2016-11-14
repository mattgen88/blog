package admin

import (
	"net/http"
  "time"
  "fmt"
  "log"
	"github.com/pmoule/go2hal/hal"
  "github.com/dgrijalva/jwt-go"

	"github.com/mattgen88/blog/util"
  "github.com/mattgen88/blog/models"
)

type AuthClaims struct {
  Username string `json:"username"`
  Role string `json:"role"`
  jwt.StandardClaims
}

// AuthHandler handles request to authenticate and will issue a JWT
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
  // @TODO look at r.Method for GET vs POST and provide some information about correctly authenticating
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

  model := models.NewSQLUser(r.FormValue("username"), h.db)
  err := model.Populate()
  if err != nil || !model.Authenticate(r.FormValue("password")) {
    w.WriteHeader(http.StatusForbidden)
    w.Write(util.JSONify(root))
    return
  }

  now := time.Now()
  expires := now.Add(time.Minute*5)

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
    root.Data()["err"] = fmt.Sprintf("%s", err)
  } else {
    root.Data()["jwt"] = tokenString
    cookie := http.Cookie{
      Name: "jwt",
      Value: tokenString,
      Secure: true,
      HttpOnly: true,
      Expires: expires,
    }
    http.SetCookie(w, &cookie)
  }
  w.Write(util.JSONify(root))

}

// AuthTest tests auth
func (h *Handler) AuthTest(w http.ResponseWriter, r *http.Request) {
  log.Println("Entered auth test")
	root := hal.NewResourceObject()

	link := &hal.LinkObject{Href: r.URL.Path}

	self := hal.NewSelfLinkRelation()
	self.SetLink(link)

	root.AddLink(self)

  root.Data()["test"] = "success"
  w.Write(util.JSONify(root))

}
