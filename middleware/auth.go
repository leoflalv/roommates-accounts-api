package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leoflalv/roommates-accounts-api/constants"
	"github.com/leoflalv/roommates-accounts-api/utils"
)

func AuthVerification(function http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer")

		// Verify issues getting cookies
		if len(splitToken) != 2 {
			utils.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		reqToken = strings.TrimSpace(splitToken[1])
		_, parseErr := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unauthorized")
			}
			return []byte(constants.JWT_SECRET_KEY), nil
		})

		if parseErr != nil {
			utils.HttpError(w, http.StatusUnauthorized, parseErr.Error())
			return
		}

		function(w, r)
	}
}
