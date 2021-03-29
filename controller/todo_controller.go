package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/model"
	"todo/service"

	"github.com/go-chi/chi/v5"
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

// Create new todo godoc
// @tags todo-manager-apis
// @Summary Create new Todo
// @Description Create new Todo
// @Accept json
// @Produce json
// @Param TodoInfo body model.Todo true "Todo information"
// @Success 200
// @Router /work/create [post]
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
	w.Write([]byte("INSERTED SUCCESSFULLY"))
}

// GetAll gets all Todos
// @tag todo-manager-apis
// @Summary Get all Todos
// @Description Get all Todos
// @Accept json
// @Produce json
// @Success 200
// @Router /work/all [get]
func (c *todoController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todos := c.todoService.GetAll()
	log.Println(todos)
	// render.JSON(w, r, todos)
	jsonResponse := struct {
		Message string
		Todos   []model.Todo
	}{
		Message: "get all successful",
		Todos:   todos,
	}
	json.NewEncoder(w).Encode(jsonResponse)

}

// GetByID gets todo by its ID
// @tag todo-manager-apis
// @Summary gets todo by its ID
// @Description gets todo by its ID
// @Accept json
// @Produce json
// @Param id path integer true "ID of needed todo"
// @Success 200
// @Router /work/search/{id} [get]
func (c *todoController) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	strID := chi.URLParam(r, "id")
	log.Println("90", strID)
	intID, cvtErr := strconv.Atoi(strID)

	if cvtErr != nil {
		log.Println("error todo_controller/getbyID<cvt>", cvtErr.Error())
		w.Write([]byte(cvtErr.Error()))
	}

	todoObject, err := c.todoService.GetByID(intID)
	jsonResponse := struct {
		Message    string
		TodoObject *model.Todo
	}{
		Message:    "Get by ID successful",
		TodoObject: todoObject,
	}
	if err != nil {
		log.Println("error encountered!/todo_controller", err.Error())
	} else {
		json.NewEncoder(w).Encode(jsonResponse)
	}

}

// Update updates information of TodoID
// @tag todo-manager-apis
// @Summary Update an ID specified Todo
// @Description Update an ID specified Todo
// @Accept json
// @Produce json
// @Param id path integer true "ID of the to be updated Todo"
// @Param UpdateContent body model.Todo true "UpdateContent information"
// @Success 200
// @Router /work/updater/{id} [put]
func (c *todoController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	strID := chi.URLParam(r, "id")
	intID, _ := strconv.Atoi(strID)

	var updateContent model.Todo
	decodeErr := json.NewDecoder(r.Body).Decode(&updateContent)
	newContent, updateErr := c.todoService.Update(intID, &updateContent)
	jsonResponse := struct {
		Message       string
		UpdateContent *model.Todo
	}{
		UpdateContent: newContent,
	}
	if decodeErr != nil {
		jsonResponse.Message = "Decode error:" + decodeErr.Error()
		json.NewEncoder(w).Encode(jsonResponse)

	}
	if updateErr != nil {
		jsonResponse.Message += "Update error:" + updateErr.Error()
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
	jsonResponse.Message = "Updated Successful"
	json.NewEncoder(w).Encode(jsonResponse)
}

// Delete deletes an ID specified Todo
// @tag todo-manager-apis
// @Summary Delete a Todo
// @Description Delete a Todo with an ID specified
// @Accept json
// @Produce json
// @Param id path integer true "ID of the to be deleted Todo"
// @Success 200
// @Router /work/delete/{id} [delete]
func (c *todoController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	strID := chi.URLParam(r, "id")
	intID, cvtErr := strconv.Atoi(strID)
	jsonResponse := struct {
		Message string
	}{
		Message: "-",
	}
	if cvtErr != nil {
		jsonResponse.Message = "convert error:" + cvtErr.Error()
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
	deleteErr := c.todoService.Delete(intID)

	if deleteErr != nil {
		jsonResponse.Message = "Delete error:" + deleteErr.Error()
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	jsonResponse.Message = "Deleted todo ID " + strconv.Itoa(intID) + " successful"
	json.NewEncoder(w).Encode(jsonResponse)

}

func NewTodoController() TodoController {
	todoService := service.NewTodoService()
	return &todoController{
		todoService: todoService,
	}
}
