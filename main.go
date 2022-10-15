package main

import (
	"fmt"

	"github.com/leoflalv/roommates-accounts-api/connection"
)

//Connection mongoDB with helper class
var db = connection.ConnectDB("roommate_accounts")

func main() {

	a := db.GetColletion("users")

	if a == nil {
		fmt.Println("Wierd")
	}
}
