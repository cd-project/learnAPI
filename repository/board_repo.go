package repository

import (
	"log"
	"strconv"
	"todo/infrastructure"
	"todo/model"
)

type boardRepository struct {
	// empty structure
}

func (b *boardRepository) CreateBoard(newBoard *model.Board, uid int) (*model.Board, error) {
	db := infrastructure.GetDB()

	err := db.Debug().Exec("INSERT INTO boards (title, description, profile_id) VALUES (?,?,?)", newBoard.Title, newBoard.Description, uid).Error
	if err != nil {
		log.Println("error REPOSITORY/CreateBoard", err.Error())
		return nil, err
	}

	return newBoard, nil
}

// update specified todo(id) of specified board(id).

func (b *boardRepository) DeleteBoard(boardID int) (error, error) {
	db := infrastructure.GetDB()

	err_boards := db.Debug().Exec("DELETE FROM boards WHERE id = ?", boardID).Error
	err_todos := db.Debug().Exec("DELETE FROM todos WHERE boardid = ?", boardID).Error
	if err_boards != nil || err_todos != nil {
		log.Println("error REPOSITORY/DeleteBoard", err_boards.Error(), err_todos.Error())
	}

	return err_boards, err_todos
}

func (b *boardRepository) GetByUserID(uid int) []model.Board {
	db := infrastructure.GetDB()

	var userBoards []model.Board
	err := db.Debug().Raw("SELECT * FROM boards WHERE profile_id = ?", uid).Scan(&userBoards).Error

	if err != nil {
		log.Println("error REPOSITORY/GetByUserID", err.Error())
		return nil
	}

	return userBoards

}

func (b *boardRepository) GetAllBoard() []model.Board {
	db := infrastructure.GetDB()

	var allBoard []model.Board
	err := db.Debug().Raw("SELECT * FROM boards").Scan(&allBoard).Error
	if err != nil {
		log.Println("error REPOSITORY/GetAllBoard", err.Error())
		return nil
	}

	log.Println("Get All Board Successfully!")
	return allBoard
}

func (b *boardRepository) UpdateBoard(boardID int, updateContent *model.Board) error {
	db := infrastructure.GetDB()

	err := db.Debug().Exec("UPDATE boards SET title = ?, description = ? WHERE id = ?", updateContent.Title, updateContent.Description, boardID).Error
	if err != nil {
		log.Println("error REPOSITORY/UpdateBoard", err.Error())
	}

	return err
}
func (b *boardRepository) FilterForSystem(filterContent *model.Board) []model.Board {
	db := infrastructure.GetDB()

	sqlString := "id > 0"
	filterUID := filterContent.ProfileID

	if filterContent.Title != "" {
		sqlString += " AND title ILIKE '%" + filterContent.Title + "%'"
	}
	if filterContent.Description != "" {
		sqlString += " AND description ILIKE '%" + filterContent.Description + "%'"
	}
	if filterUID > 0 {
		sqlString += " AND profile_id = " + strconv.Itoa(filterContent.ProfileID)
	}

	var filteredOutput []model.Board
	db.Debug().Table("boards").Where(sqlString).Scan(&filteredOutput)

	return filteredOutput

}

func NewBoardRepository() model.BoardRepository {
	return &boardRepository{}
}
