package main

import (
	"log"
)

func main() {
	port := GetConfig().Port

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":"+port, store)
	server.Run()
}
