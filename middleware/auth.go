package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leoflalv/roommates-accounts-api/constants"
	"github.com/leoflalv/roommates-accounts-api/utils"
)

func AuthVerification(function http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")

		// Verify issues getting cookies
		if err != nil {
			utils.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		_, parseErr := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unauthorized")
			}
			return []byte(constants.JWT_SECRET_KEY), nil
		})

		if parseErr != nil {
			utils.HttpError(w, http.StatusUnauthorized, err.Error())
			return
		}

		function(w, r)
	}
}
