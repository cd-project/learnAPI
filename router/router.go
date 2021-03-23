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
	boardController := controller.NewBoardController()

	// TEST
	myRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")
	// Routers for TODOS table
	myRouter.HandleFunc("/work/{boardid}/create", todoController.Create).Methods("POST")
	myRouter.HandleFunc("/work/all", todoController.GetAll).Methods("GET")
	myRouter.HandleFunc("/work/search/{id}", todoController.GetByID).Methods("GET")
	myRouter.HandleFunc("/work/updater/{id}", todoController.Update).Methods("PUT")
	myRouter.HandleFunc("/work/delete/{id}", todoController.Delete).Methods("DELETE")

	// Routers for BOARD table
	myRouter.HandleFunc("/user/{uid}/board/create", boardController.CreateBoard).Methods("POST")
	myRouter.HandleFunc("/board/{boardid}/update", boardController.UpdateBoard).Methods("PUT")
	myRouter.HandleFunc("/board/delete/{boardid}", boardController.DeleteBoard).Methods("DELETE")
	myRouter.HandleFunc("/user/{uid}/allBoard", boardController.GetByUserID).Methods("GET")
	myRouter.HandleFunc("/sys/allBoard", boardController.GetAllBoard).Methods("GET")
	myRouter.HandleFunc("/sys/filter", boardController.FilterForSystem).Methods("GET")

	// Routers for USER table
	return myRouter
}
