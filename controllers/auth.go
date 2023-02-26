package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/leoflalv/roommates-accounts-api/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserService models.UserService
}

// .
// POST register
// .
func (ac AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var resp Response[struct{}]

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp = Response[struct{}]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
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
