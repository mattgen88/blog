package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mattgen88/blog/models"
	"github.com/mattgen88/haljson"
)

// AuthMiddleware wraps something requiring auth in the form of a jwt
func AuthMiddleware(handler http.Handler, jwtKey string, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		root := haljson.NewResource()
		root.Self(r.URL.Path)

		var success bool

		// Snag JWT, verify, validate or redirect to auth endpoint
		cookie, err := r.Cookie("access.jwt")
		if err != nil {
			success = false
		} else {
			success, ctx = validateToken(ctx, cookie, jwtKey)
		}

		if !success {
			log.Println("access.jwt failed to validate")
			// Try refresh.jwt
			cookie, err := r.Cookie("refresh.jwt")
			if err != nil {
				success = false
			} else {
				success, ctx = validateToken(ctx, cookie, jwtKey)
				if success {
					log.Println("refresh.jwt validated, updating access.jwt")
					if val := ctx.Value(userDataKey("user_data")); val != nil {
						mapClaims := val.(jwt.MapClaims)
						claims := mapClaims["Claims"].(map[string]interface{})
						if username, ok := claims["username"]; ok {

							model := models.NewSQLUser(username.(string), db)
							err := model.Populate()
							if err != nil {
								log.Println(err)
								success = false
							} else {
								now := time.Now()

								accessExpires := now.Add(time.Minute * 5)

								// Create the Claims
								accessClaims := Claims{
									model.Username,
									model.Role,
									jwt.StandardClaims{
										ExpiresAt: accessExpires.Unix(),
										Issuer:    "test",
									},
								}
								token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

								ctx = context.WithValue(ctx, userDataKey("user_data"), token.Claims.(Claims))

								accessCookie, accessErr := createJwt("access.jwt", accessExpires, &accessClaims, jwtKey)
								if accessErr != nil {
									log.Println(accessErr)
									success = false
								} else {
									http.SetCookie(w, accessCookie)
								}
							}
						} else {
							log.Println("No username")
							success = false
						}
					} else {
						log.Println("No user_data")
						success = false
					}
				} else {
					log.Println("refresh.jwt failed to validate")
				}
			}

		}

		if !success {
			root.Data["error"] = "Access denied"
			w.WriteHeader(http.StatusForbidden)
			json, err := json.Marshal(root)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(json)
			return
		}
		handler.ServeHTTP(w, r.WithContext(ctx))

	})
}

func validateToken(ctx context.Context, cookie *http.Cookie, jwtKey string) (bool, context.Context) {
	success := true

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

	if err != nil {
		success = false
	}

	if !token.Valid {
		success = false
	}

	if _, ok := err.(*jwt.ValidationError); ok {
		success = false
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		success = false
	}

	ctx = context.WithValue(ctx, userDataKey("user_data"), token.Claims.(jwt.MapClaims))

	return success, ctx
}
