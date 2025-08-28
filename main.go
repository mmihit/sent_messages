package main

import (
	"fmt"
	"log"
	"net/http"

	"my_app/db"
	"my_app/handlers"

	"github.com/gorilla/mux"
)

func main() {
	Database, err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	h := &handlers.Handlers{DB: Database}

	r := mux.NewRouter()

	r.HandleFunc("/login", h.LoginHandler).Methods("POST")
	r.HandleFunc("/register", h.RegisterHandler).Methods("POST")

	serveur := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("server is running in http://localhost" + serveur.Addr)

	log.Fatal(serveur.ListenAndServe())
}
