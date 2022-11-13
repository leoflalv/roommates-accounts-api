package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/controllers"
	"github.com/leoflalv/roommates-accounts-api/routes"
	"github.com/leoflalv/roommates-accounts-api/services"
)

type App struct {
	RoutesManager routes.RoutesManager
	Initilized    bool
}

func (app *App) Initialize(dbName string) {
	db := connection.ConnectDB(dbName)

	userService := &services.UserService{Db: db}
	userController := controllers.UserController{UserService: userService}

	router := mux.NewRouter()
	routesManager := routes.RoutesManager{
		UserController: userController,
		Router:         *router,
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
	app := App{}
	app.Initialize("roommate_accounts")
	app.Run("127.0.0.1:8000")
}
