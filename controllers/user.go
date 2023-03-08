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
	var resp utils.Response[[]models.User]

	if err != nil {
		resp = utils.Response[[]models.User]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = utils.Response[[]models.User]{Data: users, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// GET user/:id
// .
func (uc UserController) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp utils.Response[models.User]

	cookie, err := r.Cookie("jwt")
	// Verify issues getting cookies
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(constants.JWT_SECRET_KEY), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := claims["issuer"].(string)

	user, err := uc.UserService.GetUserById(userId)

	if err == mongo.ErrNoDocuments {
		resp = utils.Response[models.User]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
	}

	if err != nil {
		resp = utils.Response[models.User]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp = utils.Response[models.User]{Data: *user, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// DELETE user/delete/:id
// .
func (uc UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	var resp utils.Response[string]

	err := uc.UserService.RemoveUser(id)

	if err == mongo.ErrNoDocuments {
		resp = utils.Response[string]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		resp = utils.Response[string]{Success: false, Errors: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = utils.Response[string]{Data: id, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// UPDATE user/update
// .
func (uc UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var resp utils.Response[models.User]

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp = utils.Response[models.User]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := uc.UserService.UpdateUser(&user)

	if err == mongo.ErrNoDocuments {
		resp = utils.Response[models.User]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	if err != nil {
		resp = utils.Response[models.User]{Success: false, Errors: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = utils.Response[models.User]{Data: user, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}
