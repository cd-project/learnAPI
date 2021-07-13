package middlewares

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/jwtauth"
)

func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fnc := func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			role, _ := claims["role"].(string)
			log.Println("Role: ", role)
			log.Println("Path: ", r.URL.Path)
			log.Println("Method: ", r.Method)

			result, err := e.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				log.Println("23, author ", err.Error())
				w.WriteHeader(http.StatusForbidden)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			log.Println("Line 28 authorization.go: authorization status", result)
			if result {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

		}
		return http.HandlerFunc(fnc)
	}
}
