package repository

import (
	"todo/infrastructure"
	"todo/model"
)

type todoRepository struct {
	// this file deals with database directly.
}

// Insert add new todo
func (r *todoRepository) Insert(new *model.Todo) (*model.Todo, error) {
	db := infrastructure.GetDB()

	err := db.Exec("INSERT INTO todos(title, description, finished) VALUES (?,?,?)", new.Title, new.Description, new.Finished).Error
	if err != nil {
		return nil, err
	}
	return new, nil
}

// return a slice of all Todo
func (r *todoRepository) GetAll() []model.Todo {
	// TODO: Implement
	db := infrastructure.GetDB()

	var res []model.Todo
	db.Table("todos").Select("*").Scan(&res)
	return res
}

// only viable when there is no duplicate ID?
func (r *todoRepository) GetById(id int) ([]model.Todo, error) {
	// TODO: Implement
	db := infrastructure.GetDB()

	var res []model.Todo
	err := db.Exec("SELECT * FROM todos WHERE id = ?", id).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, err
}

func (r *todoRepository) Update(id int, new *model.Todo) error {
	panic("not implemented")
}

func (r *todoRepository) Delete(id int) error {
	db := infrastructure.GetDB()
	err := db.Exec("DELETE FROM todos WHERE id = ?", id).Error

	return err

}

func NewTodoRepository() model.TodoRepository {
	return &todoRepository{}
}
