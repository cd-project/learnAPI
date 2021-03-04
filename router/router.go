package router

import (
	"net/http"
	"todo/controller"

	"github.com/gorilla/mux"
)

func Router() http.Handler {
	myRouter := mux.NewRouter().StrictSlash(true)

	// Declare controller
	todoController := controller.NewTodoController()

	// TEST
	myRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	myRouter.HandleFunc("/work", todoController.Create).Methods("POST")
	myRouter.HandleFunc("/work", todoController.GetAll).Methods("GET")
	myRouter.HandleFunc("/work/{id}", todoController.GetById).Methods("GET")
	myRouter.HandleFunc("/work/{id}", todoController.Update).Methods("PUT")
	myRouter.HandleFunc("/work/{id}", todoController.Delete).Methods("DELETE")

	return myRouter
}
