package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Duke9289/candidate-screening/structs"

	"github.com/dgrijalva/jwt-go"
)

type Exception structs.Exception

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("Authorization") //Grab the token from the header

		header = strings.TrimPrefix(header, "Bearer")
		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
			return
		}
		tk := &structs.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JwtVerifyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("Authorization") //Grab the token from the header

		header = strings.TrimPrefix(header, "Bearer")
		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
			return
		}
		tk := &structs.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}

		//&{1 tariq tariq.riahi@gmail.com 0xc000210f60}
		fmt.Println("Admin page request by: ", tk)
		if tk.Role == "basic" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "You are not allowed to be here."})
		} else if tk.Role == "admin" {
			//Check the context here
			ctx := context.WithValue(r.Context(), "user", tk)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "You are not allowed to be here."})
		}

	})
}