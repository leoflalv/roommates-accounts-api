package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leoflalv/roommates-accounts-api/constants"
)

func AuthVerification(function http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("jwt")

		_, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.SECRET_JWT_KEY), nil
		})

		if err != nil {
			resp := struct{ success bool }{success: false}
			w.WriteHeader(http.StatusUnauthorized)
			jsonResponse, _ := json.Marshal(resp)
			w.Write(jsonResponse)
		} else {
			function(w, r)
		}
	}
}
