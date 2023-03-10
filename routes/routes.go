package routes

import (
	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/controllers"
	"github.com/leoflalv/roommates-accounts-api/middleware"
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
	rm.Router.HandleFunc("/me", middleware.AuthVerification(rm.UserController.GetMe)).Methods("GET")
	rm.Router.HandleFunc("/users", middleware.AuthVerification(rm.UserController.GetUsersHandler)).Methods("GET")
	rm.Router.HandleFunc("/user/delete/{id}", middleware.AuthVerification(rm.UserController.DeleteUserHandler)).Methods("DELETE")
	rm.Router.HandleFunc("/user/update", middleware.AuthVerification(rm.UserController.UpdateUserHandler)).Methods("PUT")

	// payment_logs
	rm.Router.HandleFunc("/payment-logs", middleware.AuthVerification(rm.PaymentLogController.GetPaymentLogsHandler)).Methods("GET")
	rm.Router.HandleFunc("/payment-log/{id}", middleware.AuthVerification(rm.PaymentLogController.GetPaymentLogsByIdHandler)).Methods("GET")
	rm.Router.HandleFunc("/payment-log/create", middleware.AuthVerification(rm.PaymentLogController.CreatePaymentLogHandler)).Methods("POST")
	rm.Router.HandleFunc("/payment-log/delete/{id}", middleware.AuthVerification(rm.PaymentLogController.DeletePaymentLogHandler)).Methods("DELETE")

	initialized = true
}
