package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Persons []Person `json:"persons"`
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/server-status", serverCheck).Methods("GET")
	router.HandleFunc("/data", Data).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)
}

func serverCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/server-status' endpointt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func Data(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/data' endpoint")
	var response Response
	persons := prepareResponse()

	response.Persons = persons

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func prepareResponse() []Person {
	var persons []Person

	var person Person
	person.Id = 1
	person.FirstName = "Alex"
	person.LastName = "Walsh"
	persons = append(persons, person)

	person.Id = 2
	person.FirstName = "Daniel"
	person.LastName = "Lira"
	persons = append(persons, person)

	person.Id = 3
	person.FirstName = "Jose"
	person.LastName = "Gonzales"
	persons = append(persons, person)
	return persons
}
