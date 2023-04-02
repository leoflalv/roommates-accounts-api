package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/controllers"
	"github.com/leoflalv/roommates-accounts-api/routes"
	"github.com/leoflalv/roommates-accounts-api/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	RoutesManager routes.RoutesManager
	Db            *mongo.Database
	Initilized    bool
}

func (app *App) Initialize() {

	userService := &services.UserService{Db: app.Db}
	userController := controllers.UserController{UserService: userService}

	authController := controllers.AuthController{UserService: userService}

	paymentLogService := &services.PaymentLogService{Db: app.Db}
	paymentLogController := controllers.PaymentLogController{UserService: userService, PaymentLogService: paymentLogService}

	router := mux.NewRouter()
	routesManager := routes.RoutesManager{
		AuthController:       authController,
		UserController:       userController,
		PaymentLogController: paymentLogController,
		Router:               *router,
	}
	routesManager.Intialize()
	app.RoutesManager = routesManager
	app.Initilized = true
}

func (app *App) Run(addr string) {
	if !app.Initilized {
		log.Fatal("You have to initialize the routes manager before use it.")
	}

	http.Handle("/", &app.RoutesManager.Router)

	fmt.Printf("Running server in %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, &app.RoutesManager.Router))
}

func main() {
	settings := connection.Settings{
		MongoDBUser:     os.Getenv("MONGO_DB_USER"),
		MongoDBPassword: os.Getenv("MONGO_DB_PASSWORD"),
	}

	db := connection.ConnectDB("roommate_accounts", settings)
	app := App{Db: db}
	app.Initialize()
	app.Run(":3000")
}
