package main

import (
	"github.com/ZaharBorisenko/Banking_App_Goland/api"
)

func main() {
	server := api.NewAPIServer(":8080")
	server.Run()
}
