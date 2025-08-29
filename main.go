package main

import (
	"fmt"
	"log"
	"net/http"

	"my_app/auth"
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
	r.Handle("/add_new_clinique", auth.RequireRole("admin")(http.HandlerFunc(h.AddClinique))).Methods("POST")
	r.Handle("/add_new_surgery_date", auth.RequireRole("admin", "clinique")(http.HandlerFunc(h.AddPatient))).Methods("POST")

	serveur := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("server is running in http://localhost" + serveur.Addr)

	log.Fatal(serveur.ListenAndServe())
}
