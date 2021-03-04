package service

import (
	"todo/model"
	"todo/repository"
)

type TodoService interface {
	Create(new *model.Todo) (*model.Todo, error)
	GetAll() ([]model.Todo, error)
	GetById(id int) (*model.Todo, error)
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
	return todo, nil
}

func (s *todoService) GetAll() ([]model.Todo, error) {
	panic("not implemented") // TODO: Implement
}

func (s *todoService) GetById(id int) (*model.Todo, error) {
	panic("not implemented") // TODO: Implement
}

func (s *todoService) Update(id int, new *model.Todo) (*model.Todo, error) {
	panic("not implemented") // TODO: Implement
}

func (s *todoService) Delete(id int) error {
	panic("not implemented") // TODO: Implement
}

func NewTodoService() TodoService {
	todoRepo := repository.NewTodoRepository()
	return &todoService{
		todoRepository: todoRepo,
	}
}
