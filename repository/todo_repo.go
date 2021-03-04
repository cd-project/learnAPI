package repository

import (
	"todo/infrastructure"
	"todo/model"
)

type todoRepository struct {
}

func (r *todoRepository) Insert(new *model.Todo) (*model.Todo, error) {
	db := infrastructure.GetDB()

	err := db.Exec("INSERT INTO todos(title, description, finished) VALUES (?,?,?)", new.Title, new.Description, new.Finished).Error
	if err != nil {
		return nil, err
	}
	return new, nil
}

func (r *todoRepository) GetAll() []model.Todo {
	panic("not implemented") // TODO: Implement
}

func (r *todoRepository) GetById(id int) (*model.Todo, error) {
	panic("not implemented") // TODO: Implement
}

func (r *todoRepository) Update(id int, new *model.Todo) error {
	panic("not implemented") // TODO: Implement
}

func (r *todoRepository) Delete(id int) error {
	panic("not implemented") // TODO: Implement
}

func NewTodoRepository() model.TodoRepository {
	return &todoRepository{}
}
