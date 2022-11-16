package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/models"
)

type UserController struct {
	UserService models.UserService
}

//
// users
//
func (uc UserController) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserService.GetAllUsers()
	var resp Response[[]models.User]

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		resp = Response[[]models.User]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp = Response[[]models.User]{Data: users, Success: true}
	jsonResponse, err := json.Marshal(resp)

	w.Write(jsonResponse)
}

//
// user/:id
//
func (uc UserController) GetUsersByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var resp Response[models.User]
	user, err := uc.UserService.GetUserById(id)

	w.Header().Set("Content-Type", "application/json")

	if user == nil {
		resp = Response[models.User]{Success: false, Errors: "No documents with this id."}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		resp = Response[models.User]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp = Response[models.User]{Data: *user, Success: true}
	jsonResponse, err := json.Marshal(resp)

	w.Write(jsonResponse)
}
