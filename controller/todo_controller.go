package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/model"
	"todo/service"
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
	if err != nil {
		// bad request
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		log.Println(err)
		return
	}
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
	panic("not implemented") // TODO: Implement
}

func (c *todoController) GetById(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (c *todoController) Update(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (c *todoController) Delete(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func NewTodoController() TodoController {
	todoService := service.NewTodoService()
	return &todoController{
		todoService: todoService,
	}
}
