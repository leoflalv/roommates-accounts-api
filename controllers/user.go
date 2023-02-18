package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	UserService models.UserService
}

//
// GET users
//
func (uc UserController) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := uc.UserService.GetAllUsers()
	var resp Response[[]models.User]

	if err != nil {
		resp = Response[[]models.User]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[[]models.User]{Data: users, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

//
// GET user/:id
//
func (uc UserController) GetUsersByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	var resp Response[models.User]
	user, err := uc.UserService.GetUserById(id)

	if err == mongo.ErrNoDocuments {
		resp = Response[models.User]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
	}

	if err != nil {
		resp = Response[models.User]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp = Response[models.User]{Data: *user, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

//
// POST user/create
//
func (uc UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var resp Response[models.User]

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp = Response[models.User]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
	}

	newUser, err := uc.UserService.CreateUser(&user)

	if err != nil {
		resp = Response[models.User]{Success: false, Errors: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[models.User]{Data: newUser, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

//
// DELETE user/delete/:id
//
func (uc UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	var resp Response[string]

	err := uc.UserService.RemoveUser(id)

	if err == mongo.ErrNoDocuments {
		resp = Response[string]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		resp = Response[string]{Success: false, Errors: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[string]{Data: id, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

//
// UPDATE user/update
//
func (uc UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var resp Response[models.User]

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp = Response[models.User]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := uc.UserService.UpdateUser(&user)

	if err == mongo.ErrNoDocuments {
		resp = Response[models.User]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
		return
	}

	if err != nil {
		resp = Response[models.User]{Success: false, Errors: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[models.User]{Data: user, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}
