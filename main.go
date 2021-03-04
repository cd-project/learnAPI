package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/agenda", getAgenda).Methods("GET")
	myRouter.HandleFunc("/work/{id}/{title}", newWork).Methods("POST")
	myRouter.HandleFunc("/work/{id}", deleteWork).Methods("DELETE")
	myRouter.HandleFunc("/work/{id}/{title}", updateWork).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Print("connected")

	InitialMigration()

	handleRequests()

}
