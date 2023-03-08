package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leoflalv/roommates-accounts-api/constants"
	"github.com/leoflalv/roommates-accounts-api/models"
	"github.com/leoflalv/roommates-accounts-api/utils"
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

// .
// POST register
// .
func (ac AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	// Verify if the structure of the json is correct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.HttpError(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Verify if the username already exist
	findUser, _ := ac.UserService.GetUserByUsername(request.Username)
	if findUser != nil {
		utils.HttpError(w, http.StatusBadRequest, "This username already exist.")
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
		utils.HttpError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.HttpSuccess[struct{}](w, http.StatusOK, nil)
}

// .
// POST login
// .
func (ac AuthController) Login(w http.ResponseWriter, r *http.Request) {

	var loginInfo LoginInfo

	// Verify if the structure of the json is correct
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		utils.HttpError(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Verify if the username exist
	user, err := ac.UserService.GetUserByUsername(loginInfo.Username)
	if err == mongo.ErrNoDocuments {
		utils.HttpError(w, http.StatusNotFound, "No user with this username")
		return
	}

	// Verify if the password is right
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginInfo.Password))
	if err != nil {
		utils.HttpError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	// Create claims
	expiredDate := time.Now().Add(time.Hour)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["issuer"] = user.ID.String()
	claims["exp"] = expiredDate.Unix()
	tokenString, err := token.SignedString([]byte(constants.JWT_SECRET_KEY))

	// Verify if something internal is wrong
	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	// If everything is ok create cookie with token
	cookie := &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: expiredDate,
		Path:    "/",
		// Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	utils.HttpSuccess[struct{}](w, http.StatusOK, nil)
}
