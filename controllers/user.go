package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leoflalv/roommates-accounts-api/models"
)

type UserController struct {
	UserService models.UserService
}

func (uc UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	uc.UserService.CreateUser(&models.User{Name: "hardcore test"})
	w.WriteHeader(http.StatusOK)
}

func (uc UserController) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserService.GetAllUsers()
	w.WriteHeader(http.StatusOK)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, "Users: %v\n", users)
}
