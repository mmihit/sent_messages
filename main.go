package main

import (
	"log"

	"my_app/db"
	"my_app/handlers"

	"github.com/gorilla/mux"
)

// import "fmt"

func main() {
	Database, err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	h := &handlers.Handlers{DB: Database}

	r := mux.NewRouter()

	r.HandleFunc("/login", h.LoginHandler).Methods("POST")
}
