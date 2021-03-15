package model

type Todo struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	Finished    bool   `gorm:"column:finished"`
}

type TodoRepository interface {
	Insert(new *Todo) (*Todo, error)
	GetAll() []Todo
	GetByID(id int) (*Todo, error)
	Update(id int, new *Todo) error
	Delete(id int) error
}
