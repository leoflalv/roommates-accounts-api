package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leoflalv/roommates-accounts-api/constants"
)

func GetIssuer(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	reqToken = strings.TrimSpace(splitToken[1])

	// Verify issues getting cookies
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(constants.JWT_SECRET_KEY), nil
	})

	if err != nil {
		return "", err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := claims["issuer"].(string)

	return userId, nil
}
