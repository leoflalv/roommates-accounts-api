package main

import (
	"fmt"

	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/services"
)

func main() {
	db := connection.ConnectDB("roommate_accounts")

	userService := services.UserService{Db: db}

	user, err := userService.GetUserById("634c15ac748267b3a765af3e")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(user)
	}

}
