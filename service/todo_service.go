package service

import (
	"log"
	"todo/model"
	"todo/repository"
)

// interface TodoService
type TodoService interface {
	Create(new *model.Todo, boardID int) (*model.Todo, error)
	GetAll() []model.Todo
	GetByID(id int) (*model.Todo, error)
	Update(id int, new *model.Todo) (*model.Todo, error)
	Delete(id int) error
}

type todoService struct {
	todoRepository model.TodoRepository
}

func (s *todoService) Create(new *model.Todo, boardID int) (*model.Todo, error) {
	todo, err := s.todoRepository.Insert(new, boardID)
	if err != nil {
		return nil, err
	}
	return todo, err
}

func (s *todoService) GetAll() []model.Todo {
	list := s.todoRepository.GetAll()

	return list
}

func (s *todoService) GetByID(id int) (*model.Todo, error) {
	todoObject, err := s.todoRepository.GetByID(id)
	failureObject := model.Todo{ID: 4444, Title: "failureObj", Description: "failureObj", Finished: false}
	if err != nil {
		log.Println("error encountered / todo_service!", err.Error())
		return &failureObject, err
	}

	return todoObject, err
}

func (s *todoService) Update(id int, new *model.Todo) (*model.Todo, error) {
	err := s.todoRepository.Update(id, new)
	if err != nil {
		log.Println(err.Error())
		return new, err
	}

	return new, err
}

func (s *todoService) Delete(id int) error {
	err := s.todoRepository.Delete(id)
	if err != nil {
		return err
	}

	return err
}

func NewTodoService() TodoService {
	todoRepo := repository.NewTodoRepository()
	x0 := &todoService{
		todoRepository: todoRepo,
	}
	return x0
}
