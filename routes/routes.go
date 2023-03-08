package routes

import (
	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/controllers"
)

type RoutesManager struct {
	Router               mux.Router
	AuthController       controllers.AuthController
	UserController       controllers.UserController
	PaymentLogController controllers.PaymentLogController
}

var initialized bool = false

func (rm *RoutesManager) Intialize() {

	//auth
	rm.Router.HandleFunc("/register", rm.AuthController.Register).Methods("POST")
	rm.Router.HandleFunc("/login", rm.AuthController.Login).Methods("POST")

	// users
	rm.Router.HandleFunc("/me", rm.UserController.GetMe).Methods("GET")
	rm.Router.HandleFunc("/users", rm.UserController.GetUsersHandler).Methods("GET")
	rm.Router.HandleFunc("/user/delete/{id}", rm.UserController.DeleteUserHandler).Methods("DELETE")
	rm.Router.HandleFunc("/user/update", rm.UserController.UpdateUserHandler).Methods("PUT")

	// payment_logs
	rm.Router.HandleFunc("/payment-logs", rm.PaymentLogController.GetPaymentLogsHandler).Methods("GET")
	rm.Router.HandleFunc("/payment-log/{id}", rm.PaymentLogController.GetPaymentLogsByIdHandler).Methods("GET")
	rm.Router.HandleFunc("/payment-log/create", rm.PaymentLogController.CreatePaymentLogHandler).Methods("POST")
	rm.Router.HandleFunc("/payment-log/delete/{id}", rm.PaymentLogController.DeletePaymentLogHandler).Methods("DELETE")

	initialized = true
}
