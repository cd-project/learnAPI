package service

import (
	"todo/model"
	"todo/repository"
)

type TodoService interface {
	Create(new *model.Todo) (*model.Todo, error)
	GetAll() []model.Todo
	GetById(id int) ([]model.Todo, error)
	Update(id int, new *model.Todo) (*model.Todo, error)
	Delete(id int) error
}

type todoService struct {
	todoRepository model.TodoRepository
}

func (s *todoService) Create(new *model.Todo) (*model.Todo, error) {
	todo, err := s.todoRepository.Insert(new)
	if err != nil {
		return nil, err
	}
	return todo, err
}

func (s *todoService) GetAll() []model.Todo {
	list := s.todoRepository.GetAll()

	return list
}

func (s *todoService) GetById(id int) ([]model.Todo, error) {
	list, err := s.todoRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (s *todoService) Update(id int, new *model.Todo) (*model.Todo, error) {
	err := s.todoRepository.Update(new.ID, new)
	if err != nil {
		return nil, err
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
