package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/model"
	"todo/service"

	"github.com/go-chi/chi/v5"
)

type BoardController interface {
	CreateBoard(w http.ResponseWriter, r *http.Request)
	UpdateBoard(w http.ResponseWriter, r *http.Request)
	DeleteBoard(w http.ResponseWriter, r *http.Request)
	GetByUserID(w http.ResponseWriter, r *http.Request)
	GetAllBoard(w http.ResponseWriter, r *http.Request)
	Filter(w http.ResponseWriter, r *http.Request)
}

type boardController struct {
	boardService service.BoardService
}

// CreateBoard create new Board godoc
// @tags board-manager-apis
// @Summary create new Board with given model
// @Description create a new board with given model
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uid path integer true "Owner of this board"
// @Param BoardInfo body model.Board true "Board information"
// @Success 200
// @Router /board/{uid}/create [post]
func (c *boardController) CreateBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := struct {
		Message string
		Data    *model.Board
	}{}
	var data model.Board
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&data)
	if decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse.Message = decodeErr.Error()
		log.Println(decodeErr)
		return
	}
	// get user id from URL
	strUID := chi.URLParam(r, "uid")
	uid, _ := strconv.Atoi(strUID)
	// create new Board
	newBoard, dbErr := c.boardService.CreateBoard(&data, uid)
	if dbErr != nil {
		log.Println("error CONTROLLER/CreateBoard", dbErr.Error())
		jsonResponse.Message += dbErr.Error()
		return
	}
	jsonResponse.Message = "Created board successful, board data:"
	jsonResponse.Data = newBoard
	jsonResponse.Data.ProfileID = uid
	json.NewEncoder(w).Encode(jsonResponse)
}

// UpdateBoard updates board of specified ID with new data.
// @tags board-manager-apis
// @Summary get board updated with new data
// @Description given new data and id, update board
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param boardid path integer true "ID of the to be updated board"
// @Param UpdateContent body model.Board true "Update content"
// @Success 200
// @Router /board/{boardid}/update [put]
func (c *boardController) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	strBID := chi.URLParam(r, "boardid")
	boardID, _ := strconv.Atoi(strBID)
	jsonResponse := struct {
		Message       string
		UpdateContent model.Board
	}{}
	var updateContent model.Board
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&updateContent)
	if decodeErr != nil {
		log.Println("error CONTROLLER/UpdateBoard", decodeErr.Error())
		jsonResponse.Message = decodeErr.Error()
		return
	} else {
		dbErr := c.boardService.UpdateBoard(boardID, &updateContent)
		if dbErr != nil {
			log.Println("error CONTROLLER/UpdateBoard", dbErr.Error())
			jsonResponse.Message += dbErr.Error()
			return
		}
	}
	jsonResponse.Message = "Updated content of board " + strBID
	jsonResponse.UpdateContent = updateContent
	json.NewEncoder(w).Encode(jsonResponse)
}

// DeleteBoard deletes board with boardID
// @tags board-manager-apis
// @Summary Board with boardID will be deleted
// @Description Board with boardID will be deleted
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param boardid path integer true "ID of the to be deleted board"
// @Success 200
// @Router /board/delete/{boardid} [delete]
func (c *boardController) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	str_boardID := chi.URLParam(r, "boardid")
	boardID, _ := strconv.Atoi(str_boardID)

	err1, err2 := c.boardService.DeleteBoard(boardID)
	if err1 != nil || err2 != nil {
		log.Println("error CONTROLLER/DeleteBoard", err1.Error(), err2.Error())
		errString := err1.Error() + ", " + err2.Error()
		w.Write([]byte(errString))
	}
}

// GetByUserID gets all board belong to an User
// @tags board-manager-apis
// @Summary gets all board belong to UserID
// @Description gets all board belong to UserID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uid path integer true "User ID"
// @Success 200
// @Router /board/{uid}/allBoard [get]
func (c *boardController) GetByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	strUID := chi.URLParam(r, "uid")
	uid, _ := strconv.Atoi(strUID)

	userBoards := c.boardService.GetByUserID(uid)

	jsonResponse := struct {
		Message string
		Boards  []model.Board
	}{
		Message: "Getting boards of user " + strUID,
		Boards:  userBoards,
	}
	json.NewEncoder(w).Encode(jsonResponse)

}

// GetAllBoard gets all boards currently available in database
// @tags board-manager-apis
// @Summary get all boards
// @Description get all boards
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200
// @Router /sys/allBoard [get]
func (c *boardController) GetAllBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allBoards := c.boardService.GetAllBoard()
	jsonResponse := struct {
		Message string
		Boards  []model.Board
	}{
		Message: "Get all boards: Successful",
		Boards:  allBoards,
	}

	json.NewEncoder(w).Encode(jsonResponse)
}

// Filter filters data with given model
// @tags board-manager-apis
// @Summary filtered data will be shown
// @Description board db will be filtered using given model
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param FilterContent body model.Board true "Filter Content"
// @Success 200
// @Router /sys/filter [put]
func (c *boardController) Filter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := struct {
		Message         string
		FilteredContent []model.Board
	}{}
	var filteredContent model.Board
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&filteredContent)
	if decodeErr != nil {
		log.Println("error CONTROLLER/Filter", decodeErr.Error())
		jsonResponse.Message = decodeErr.Error()
		return
	} else {
		filteredOutput := c.boardService.Filter(&filteredContent)
		jsonResponse.Message = "Filtered content:"
		jsonResponse.FilteredContent = filteredOutput
	}
	json.NewEncoder(w).Encode(jsonResponse)
}

func NewBoardController() BoardController {
	boardService := service.NewBoardService()

	return &boardController{
		boardService: boardService,
	}
}
