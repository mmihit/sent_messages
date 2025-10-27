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

	err = h.SchedulerTask()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/login", h.LoginHandler).Methods("POST")
	r.Handle("/add_new_clinique", auth.RequireRole("admin")(http.HandlerFunc(h.AddClinique))).Methods("POST")
	r.Handle("/add_patient", auth.RequireRole("clinique")(http.HandlerFunc(h.AddPatient))).Methods("POST")
	r.Handle("/delete_patient/{id}", auth.RequireRole("clinique")(http.HandlerFunc(h.DeletePatient))).Methods("DELETE")
	r.Handle("/get_patients_by_removal_date/{date}", auth.RequireRole("clinique")(http.HandlerFunc(h.GetPatientsByRemovalDate))).Methods("GET")
	r.Handle("/get_patient_by_id/{id}", auth.RequireRole("clinique")(http.HandlerFunc(h.GetPatientById))).Methods("GET")
	r.Handle("/get_all_patients", auth.RequireRole("clinique")(http.HandlerFunc(h.GetAllPatients))).Methods("GET")
	r.Handle("/get_patients_numbers", auth.RequireRole("clinique")(http.HandlerFunc(h.GetPatientsCount))).Methods("GET")

	serveur := &http.Server{
		Addr:    ":8080",
		Handler: h.HandleCORS(r),
	}

	fmt.Println("server is running in http://localhost" + serveur.Addr)

	log.Fatal(serveur.ListenAndServe())
}
