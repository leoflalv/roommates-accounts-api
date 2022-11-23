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

	rm.Router.HandleFunc("/users", rm.UserController.GetUsersHandler).Methods("GET")
	rm.Router.HandleFunc("/user/{id}", rm.UserController.GetUsersByIdHandler).Methods("GET")
	rm.Router.HandleFunc("/user/create", rm.UserController.CreateUserHandler).Methods("POST")
	rm.Router.HandleFunc("/user/delete/{id}", rm.UserController.DeleteUserHandler).Methods("DELETE")
	rm.Router.HandleFunc("/user/update", rm.UserController.UpdateUserHandler).Methods("PUT")

	initialized = true
}
