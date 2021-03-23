package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo/model"
	"todo/service"

	"github.com/gorilla/mux"
)

type TodoController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type todoController struct {
	todoService service.TodoService
}

func (c *todoController) Create(w http.ResponseWriter, r *http.Request) {
	// get body request and decode
	var data model.Todo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		// bad request
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	strID := vars["boardid"]
	boardID, _ := strconv.Atoi(strID)
	// create new Todo
	new, err := c.todoService.Create(&data, boardID)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	// success
	log.Println(new)
	w.Write([]byte("INSERTED SUCCESSFULLY"))
}

func (c *todoController) GetAll(w http.ResponseWriter, r *http.Request) {
	todos := c.todoService.GetAll()
	json.NewEncoder(w).Encode(todos)

	w.Write([]byte("GET ALL : SUCCESSFUL!"))
}

func (c *todoController) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	intID, _ := strconv.Atoi(strID)

	todoObject, err := c.todoService.GetByID(intID)
	if err != nil {
		log.Println("error encountered! / todo_controller", err.Error())
	} else {
		w.Write([]byte("Object Catched!"))
	}

	json.NewEncoder(w).Encode(todoObject)
}

func (c *todoController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	intID, _ := strconv.Atoi(strID)

	//get newTodo object from request body
	var newTodo model.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	fmt.Println(newTodo)
	if err != nil {
		// bad request
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		log.Println(err)
		return
	} else {
		w.Write([]byte("Updated Successfully"))
	}
	c.todoService.Update(intID, &newTodo)
}

func (c *todoController) Delete(w http.ResponseWriter, r *http.Request) {
	// get "id" from URL
	vars := mux.Vars(r)
	strID := vars["id"]
	intID, _ := strconv.Atoi(strID)

	err := c.todoService.Delete(intID)
	if err != nil {
		panic(err.Error())
	} else {
		w.Write([]byte("Deleted Successfully!"))
	}
}

func NewTodoController() TodoController {
	todoService := service.NewTodoService()
	return &todoController{
		todoService: todoService,
	}
}
