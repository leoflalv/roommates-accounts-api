package routes

import (
	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/controllers"
)

type RoutesManager struct {
	Router         mux.Router
	UserController controllers.UserController
}

var initialized bool = false

func (rm *RoutesManager) Intialize() {

	rm.Router.HandleFunc("/users", rm.UserController.GetUsersHandler).Methods("GET").Schemes("http")
	rm.Router.HandleFunc("/user/{id}", rm.UserController.GetUsersByIdHandler).Methods("GET").Schemes("http")
	// rm.Router.HandleFunc("/test/", rm.UserController.CreateUserHandler).Methods("POST")

	initialized = true
}
