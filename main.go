package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Patient struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Age           int    `json:"age"`
	Disease       string `json:"disease"`
	AdmissionDate string `json:"admission_date"`
	Status        string `json:"status"`
}

var patients []Patient
var nextID int = 1

// Create a new patient
func createPatient(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	json.NewDecoder(r.Body).Decode(&patient)
	patient.ID = nextID
	nextID++
	patients = append(patients, patient)
	json.NewEncoder(w).Encode(patient)
}

// Get all patients
func getPatients(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(patients)
}

// Get a single patient by ID
func getPatientByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, patient := range patients {
		if patient.ID == id {
			json.NewEncoder(w).Encode(patient)
			return
		}
	}
	http.Error(w, "Patient not found", http.StatusNotFound)
}

// Update a patient by ID
func updatePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, patient := range patients {
		if patient.ID == id {
			json.NewDecoder(r.Body).Decode(&patient)
			patients[i] = patient
			json.NewEncoder(w).Encode(patient)
			return
		}
	}
	http.Error(w, "Patient not found", http.StatusNotFound)
}

// Delete a patient by ID
func deletePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, patient := range patients {
		if patient.ID == id {
			patients = append(patients[:i], patients[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Patient not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/patients", createPatient).Methods("POST")
	router.HandleFunc("/patients", getPatients).Methods("GET")
	router.HandleFunc("/patients/{id}", getPatientByID).Methods("GET")
	router.HandleFunc("/patients/{id}", updatePatient).Methods("PUT")
	router.HandleFunc("/patients/{id}", deletePatient).Methods("DELETE")

	log.Println("Server running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", router))
}
