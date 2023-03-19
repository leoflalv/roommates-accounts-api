package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/constants"
	"github.com/leoflalv/roommates-accounts-api/models"
	"github.com/leoflalv/roommates-accounts-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	UserService models.UserService
}

// .
// GET users
// .
func (uc UserController) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := uc.UserService.GetAllUsers()

	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.HttpSuccess(w, http.StatusOK, &users)
}

// .
// GET user/:id
// .
func (uc UserController) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, _ := r.Cookie("jwt")

	// Verify issues getting cookies
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(constants.JWT_SECRET_KEY), nil
	})

	if err != nil {
		utils.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := claims["issuer"].(string)

	user, err := uc.UserService.GetUserById(userId)

	if err == mongo.ErrNoDocuments {
		utils.HttpError(w, http.StatusNotFound, "No documents with this id")
		return
	}

	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.HttpSuccess(w, http.StatusOK, &user)
}

// .
// DELETE user/delete/:id
// .
func (uc UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	err := uc.UserService.RemoveUser(id)

	if err == mongo.ErrNoDocuments {
		utils.HttpError(w, http.StatusNotFound, "No user with this id")
		return
	}

	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.HttpSuccess(w, http.StatusOK, &id)
}

// .
// UPDATE user/update
// .
func (uc UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := uc.UserService.UpdateUser(&user)

	if err == mongo.ErrNoDocuments {
		utils.HttpError(w, http.StatusNotFound, "No user with this id")
		return
	}

	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.HttpSuccess(w, http.StatusOK, &user)
}
