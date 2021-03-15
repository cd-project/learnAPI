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

	myRouter.HandleFunc("/work/create", todoController.Create).Methods("POST")
	myRouter.HandleFunc("/work/all", todoController.GetAll).Methods("GET")
	myRouter.HandleFunc("/work/search/{id}", todoController.GetByID).Methods("GET")
	myRouter.HandleFunc("/work/updater/{id}", todoController.Update).Methods("PUT")
	myRouter.HandleFunc("/work/delete/{id}", todoController.Delete).Methods("DELETE")

	return myRouter
}
