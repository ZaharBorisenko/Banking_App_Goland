package main

import (
	"github.com/ZaharBorisenko/Banking_App_Goland/api"
	"github.com/ZaharBorisenko/Banking_App_Goland/storage"
	"log"
)

func main() {
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", store)
	server.Run()
}
