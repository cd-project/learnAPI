package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/model"
	"todo/service"

	"github.com/gorilla/mux"
)

type BoardController interface {
	CreateBoard(w http.ResponseWriter, r *http.Request)
	UpdateBoard(w http.ResponseWriter, r *http.Request)
	DeleteBoard(w http.ResponseWriter, r *http.Request)
	GetByUserID(w http.ResponseWriter, r *http.Request)
	GetAllBoard(w http.ResponseWriter, r *http.Request)
	FilterForSystem(w http.ResponseWriter, r *http.Request)
}

type boardController struct {
	boardService service.BoardService
}

func (c *boardController) CreateBoard(w http.ResponseWriter, r *http.Request) {
	// decode request body
	var data model.Board
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		log.Println(err)
		return
	}
	log.Println(data)
	// get user id from URL
	vars := mux.Vars(r)
	strID := vars["uid"]
	uid, _ := strconv.Atoi(strID)
	// create new Board
	newBoard, err := c.boardService.CreateBoard(&data, uid)
	if err != nil {
		log.Println("error CONTROLLER/CreateBoard", err.Error())
		w.Write([]byte(err.Error()))
	}

	log.Println(newBoard)
	w.Write([]byte("CREATED SUCCESSFULLY"))
}

func (c *boardController) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	var data model.Board
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		// bad request
		log.Println("error CONTROLLER/UpdateBoard", err.Error())
	} else {
		c.boardService.UpdateBoard(data.ID, &data)
	}
}

func (c *boardController) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str_boardID := vars["boardid"]
	boardID, _ := strconv.Atoi(str_boardID)

	err1, err2 := c.boardService.DeleteBoard(boardID)
	if err1 != nil || err2 != nil {
		log.Println("error CONTROLLER/DeleteBoard", err1.Error(), err2.Error())
	}
}

func (c *boardController) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strUID := vars["uid"]
	uid, _ := strconv.Atoi(strUID)

	userBoards := c.boardService.GetByUserID(uid)
	json.NewEncoder(w).Encode(userBoards)

	w.Write([]byte("GetByUserID successfully!"))
}

func (c *boardController) GetAllBoard(w http.ResponseWriter, r *http.Request) {
	allBoard := c.boardService.GetAllBoard()
	json.NewEncoder(w).Encode(allBoard)

	w.Write([]byte("GET ALL BOARD: SUCCESSFUL!"))
}

func (c *boardController) FilterForSystem(w http.ResponseWriter, r *http.Request) {
	var filteredContent model.Board
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&filteredContent)
	if err != nil {
		log.Println("error CONTROLLER/FilterForSystem", err.Error())
	} else {
		filteredOutput := c.boardService.FilterForSystem(&filteredContent)
		json.NewEncoder(w).Encode(filteredOutput)
	}
}

func NewBoardController() BoardController {
	bsv := service.NewBoardService()

	return &boardController{
		boardService: bsv,
	}
}
