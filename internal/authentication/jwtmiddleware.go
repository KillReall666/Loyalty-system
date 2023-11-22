package jwtmiddleware

import (
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/authentication"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
)

func JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		jwtCheck := func(w http.ResponseWriter, r *http.Request) {
			extractor := request.AuthorizationHeaderExtractor
			extToken, err := extractor.ExtractToken(r)
			if err != nil {
				fmt.Println("err when extract token in jwt middleware: ", err)
			}
			x, err := authentication.GetUserID(extToken)
			if x == "" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(jwtCheck)
	}
}
