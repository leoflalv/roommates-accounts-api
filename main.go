package main

import (
	"fmt"
	"log"

	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/services"
)

func main() {
	db := connection.ConnectDB("roommate_accounts")

	userService := services.UserService{Db: db}

	user, err := userService.RemoveUser("635d87d5a2cc82dc7773059f")

	if err != nil {
		fmt.Println(">>>>")
		log.Fatal(err.Error())
	} else {
		fmt.Println(user)
	}

}
