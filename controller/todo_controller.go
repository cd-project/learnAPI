package controller

import (
	"encoding/json"
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
	GetById(w http.ResponseWriter, r *http.Request)
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
	// if err != nil {
	// 	// bad request
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	http.Error(w, http.StatusText(400), 400)
	// 	log.Println(err)
	// 	return
	// }
	// create new Todo
	new, err := c.todoService.Create(&data)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	// success
	log.Println(new)
	w.Write([]byte("INSERT SUCCESSFUL"))
}

func (c *todoController) GetAll(w http.ResponseWriter, r *http.Request) {
	todos := c.todoService.GetAll()
	json.NewEncoder(w).Encode(todos)
	log.Println(todos)
}

func (c *todoController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	intID, _ := strconv.Atoi(strID)

	list, err := c.todoService.GetById(intID)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(list)
}

func (c *todoController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	intID, _ := strconv.Atoi(strID)

	//get newTodo object from request body
	var newTodo model.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		// bad request
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		log.Println(err)
		return
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
	}
}

func NewTodoController() TodoController {
	todoService := service.NewTodoService()
	return &todoController{
		todoService: todoService,
	}
}
