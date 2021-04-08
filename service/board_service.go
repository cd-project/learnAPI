package service

import (
	"log"
	"todo/model"
	"todo/repository"
)

type BoardService interface {
	CreateBoard(newBoard *model.Board, uid int) (*model.Board, error)
	DeleteBoard(boardID int) (error, error)
	GetByUserID(uid int) []model.Board
	GetAllBoard() []model.Board
	UpdateBoard(boardID int, updateContent *model.Board) error
	Filter(filterContent *model.Board) []model.Board
}

type boardService struct {
	boardRepository model.BoardRepository
}

func (s *boardService) CreateBoard(newBoard *model.Board, uid int) (*model.Board, error) {
	board, err := s.boardRepository.CreateBoard(newBoard, uid)
	if err != nil {
		log.Println("error SERVICE/CreateBoard", err.Error())
		return nil, err
	}

	return board, nil

}

func (s *boardService) UpdateBoard(boardID int, updateContent *model.Board) error {
	err := s.boardRepository.UpdateBoard(boardID, updateContent)
	if err != nil {
		log.Println("error SERVICE/UpdateBoard", err.Error())
	}

	return err
}

func (s *boardService) DeleteBoard(boardID int) (error, error) {
	log.Println(boardID)
	err1, err2 := s.boardRepository.DeleteBoard(boardID)
	if err1 != nil || err2 != nil {
		log.Println("error SERVICE/DeleteBoard", err1.Error(), err2.Error())
	}

	return err1, err2
}

func (s *boardService) GetByUserID(uid int) []model.Board {
	userBoards := s.boardRepository.GetByUserID(uid)

	if userBoards == nil {
		log.Println("This user owns no board.")
	}

	return userBoards
}

func (s *boardService) GetAllBoard() []model.Board {
	allBoard := s.boardRepository.GetAllBoard()

	return allBoard
}

func (s *boardService) Filter(filterContent *model.Board) []model.Board {
	filteredOutput := s.boardRepository.Filter(filterContent)
	return filteredOutput
}

func NewBoardService() BoardService {
	boardRepo := repository.NewBoardRepository()
	x0 := &boardService{
		boardRepository: boardRepo,
	}

	return x0
}
