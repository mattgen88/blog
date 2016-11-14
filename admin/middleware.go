package admin

import (
  "log"
  "fmt"
	"net/http"
  "github.com/pmoule/go2hal/hal"
  "github.com/dgrijalva/jwt-go"
  	"github.com/mattgen88/blog/util"
)

func AuthMiddleware(handler http.Handler, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    root := hal.NewResourceObject()

    link := &hal.LinkObject{Href: r.URL.Path}

    self := hal.NewSelfLinkRelation()
    self.SetLink(link)

    root.AddLink(self)

		// Snag JWT, verify, validate or redirect to auth endpoint
    cookie, err := r.Cookie("jwt")
    if err != nil {
      log.Println("Failed to get cookie jtw: " + fmt.Sprintf("%s", err))

      root.Data()["error"] = "Auth not found"
      w.WriteHeader(http.StatusForbidden)
      w.Write(util.JSONify(root))
      return
    }

    // Parse takes the token string and a function for looking up the key. The latter is especially
    // useful if you use multiple keys for your application.  The standard is to use 'kid' in the
    // head of the token to identify which key to use, but the parsed token (head and claims) is provided
    // to the callback, providing flexibility.
    token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtKey), nil
    })

    if token.Valid {
        log.Println("You look nice today")
    } else if ve, ok := err.(*jwt.ValidationError); ok {
        if ve.Errors&jwt.ValidationErrorMalformed != 0 {
            log.Println("That's not even a token")
        } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
            // Token is either expired or not active yet
            log.Println("Timing is everything")
        } else {
            log.Println("Couldn't handle this token:", err)
        }
    } else {
        log.Println("Couldn't handle this token:", err)
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        log.Println(claims["role"], claims["username"], claims["nbf"])
    } else {
        root.Data()["error"] = "claims not OK"
        root.Data()["claims"] = claims
        w.WriteHeader(http.StatusForbidden)
        w.Write(util.JSONify(root))
        return
    }

		handler.ServeHTTP(w, r)
	})
}
