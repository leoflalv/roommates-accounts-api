package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leoflalv/roommates-accounts-api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Cookie struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthController struct {
	UserService models.UserService
}

var SecretJWTKey = os.Getenv("JWT_SECRET_KEY")

// .
// POST register
// .
func (ac AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response[struct{}]

	var request struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	// Verify if the structure of the json is correct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	findUser, _ := ac.UserService.GetUserByUsername(request.Username)

	// Verify if the username already exist
	if findUser != nil {
		resp = Response[struct{}]{Success: false, Errors: "This username already exist."}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	// Create the user
	user := models.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Password:  []byte(request.Password),
	}
	user.Password, _ = bcrypt.GenerateFromPassword(user.Password, 14)
	_, err := ac.UserService.CreateUser(&user)

	// Verify if everything createing the user is correct
	if err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[struct{}]{Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// POST login
// .
func (ac AuthController) Login(w http.ResponseWriter, r *http.Request) {

	var resp Response[struct{}]
	var loginInfo LoginInfo

	// Verify if the structure of the json is correct
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	// Verify if the username exist
	user, err := ac.UserService.GetUserByUsername(loginInfo.Username)
	if err == mongo.ErrNoDocuments {
		resp = Response[struct{}]{Success: false, Errors: "No user with this username"}
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	// Verify if the password is right
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginInfo.Password))
	if err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Incorrect password"}
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	// Create claims
	expiredDate := time.Now().Add(time.Hour)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Unix(expiredDate.Unix(), 0)),
	})

	token, err := claims.SignedString([]byte(SecretJWTKey))

	// Verify if something internal is wrong
	if err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Could not login"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[struct{}]{Success: true}

		// If everything is ok create cookie with token
		cookie := &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  expiredDate,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}
