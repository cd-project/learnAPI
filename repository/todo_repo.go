package repository

import (
	"log"
	"todo/infrastructure"
	"todo/model"
)

type todoRepository struct {
	// this file deals with database directly.
}

// Insert add new todo into board with boardID
func (r *todoRepository) Insert(new *model.Todo, boardID int) (*model.Todo, error) {
	db := infrastructure.GetDB()

	err := db.Debug().Exec("INSERT INTO todos(title, description, finished, boardid) VALUES (?,?,?,?)", new.Title, new.Description, new.Finished, boardID).Error
	if err != nil {
		log.Println("error", err.Error())
		return nil, err
	}
	log.Println(new)
	return new, nil
}

// return a slice of all Todo
func (r *todoRepository) GetAll() []model.Todo {
	db := infrastructure.GetDB()

	var res []model.Todo
	err := db.Table("todos").Select("*").Scan(&res).Error
	if err != nil {
		log.Println("Error encountered!", err.Error())
		return nil
	}
	log.Println("Status OK")
	return res
}

func (r *todoRepository) GetByID(id int) (*model.Todo, error) {
	db := infrastructure.GetDB()

	var res model.Todo
	err := db.Debug().Raw("SELECT * FROM todos WHERE id = ?", id).Scan(&res).Error
	log.Println(res)
	if err != nil {
		log.Println("error encountered! todo_repo", err.Error())
		return nil, err
	}

	return &res, err
}

func (r *todoRepository) UpdateTodoInBoard(boardID int, todoID int, newTodo *model.Todo) (*model.Todo, error) {
	db := infrastructure.GetDB()

	err := db.Debug().Exec("UPDATE todos SET title = ?, description = ?, finished = ? WHERE boardid = ? AND id = ?", newTodo.Title, newTodo.Description, newTodo.Finished, boardID, todoID).Error
	if err != nil {
		log.Println("error REPOSITORY/UpdateTodoInBoard", err.Error())
		return newTodo, err
	}

	return newTodo, err
}
func (r *todoRepository) Update(id int, new *model.Todo) error {
	db := infrastructure.GetDB()

	err := db.Debug().Exec("UPDATE todos SET title = ?, description = ?, finished = ? WHERE id = ?", new.Title, new.Description, new.Finished, id).Error
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(new)

	return err
}

func (r *todoRepository) Delete(id int) error {
	db := infrastructure.GetDB()
	err := db.Debug().Exec("DELETE FROM todos WHERE id = ?", id).Error

	return err

}

func NewTodoRepository() model.TodoRepository {
	return &todoRepository{}
}
