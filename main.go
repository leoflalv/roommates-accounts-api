package main

import (
	"fmt"
	"log"

	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/services"
)

func main() {
	db := connection.ConnectDB("roommate_accounts")

	// userService := services.UserService{Db: db}
	paymentLogService := services.PaymentLogService{Db: db}

	// user, err := userService.RemoveUser("635d87d5a2cc82dc7773059f")
	paymentLogs, err := paymentLogService.GetPaymentsLogsByPayer("634c15ac748267b3a765af3e")

	if err != nil {
		fmt.Println(">>>>")
		log.Fatal(err.Error())
	} else {
		fmt.Println(paymentLogs)
	}

}
