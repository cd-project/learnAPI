package model

type Board struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	ProfileID   int    `gorm:"column:profile_id"`
}

// Each user has an unique set of Boards.
type BoardRepository interface {
	CreateBoard(newBoard *Board, uid int) (*Board, error)
	DeleteBoard(boardID int) (error, error)
	GetByUserID(uid int) []Board
	GetAllBoard() []Board
	UpdateBoard(boardID int, updateContent *Board) error
	Filter(filterContent *Board) []Board
}
