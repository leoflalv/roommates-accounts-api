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

type AuthResponse struct {
	Token string `json:"token"`
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

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	findUser, _ := ac.UserService.GetUserByUsername(request.Username)

	if findUser != nil {
		resp = Response[struct{}]{Success: false, Errors: "This username already exist."}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	user := models.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Password:  []byte(request.Password),
	}

	user.Password, _ = bcrypt.GenerateFromPassword(user.Password, 14)

	_, err := ac.UserService.CreateUser(&user)

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
	w.Header().Set("Content-Type", "application/json")

	var resp Response[AuthResponse]
	var loginInfo LoginInfo

	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		resp = Response[AuthResponse]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	user, err := ac.UserService.GetUserByUsername(loginInfo.Username)

	if err == mongo.ErrNoDocuments {
		resp = Response[AuthResponse]{Success: false, Errors: "No user with this username"}
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginInfo.Password))

	if err != nil {
		resp = Response[AuthResponse]{Success: false, Errors: "Incorrect password"}
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	} else {
		resp = Response[AuthResponse]{Success: true}
		w.WriteHeader(http.StatusOK)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Add(time.Hour*24).Unix(), 0)),
	})

	token, err := claims.SignedString([]byte(SecretJWTKey))

	if err != nil {
		resp = Response[AuthResponse]{Success: false, Errors: "Could not login"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[AuthResponse]{Success: true, Data: AuthResponse{Token: token}}
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}
