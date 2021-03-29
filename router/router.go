package router

import (
	"log"
	"net/http"
	"os"
	"todo/controller"

	_ "todo/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errLog  = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Router() http.Handler {
	myRouter := chi.NewRouter()

	// Declare controller
	todoController := controller.NewTodoController()
	boardController := controller.NewBoardController()

	// TEST
	myRouter.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Routers for TODOS table
	myRouter.Post("/work/create", todoController.Create)
	myRouter.Get("/work/all", todoController.GetAll)
	myRouter.Get("/work/search/{id}", todoController.GetByID)
	myRouter.Put("/work/updater/{id}", todoController.Update)
	myRouter.Delete("/work/delete/{id}", todoController.Delete)

	// Routers for BOARD table
	myRouter.Post("/user/{uid}/board/create", boardController.CreateBoard)
	myRouter.Put("/board/{boardid}/update", boardController.UpdateBoard)
	myRouter.Delete("/board/delete/{boardid}", boardController.DeleteBoard)
	myRouter.Get("/user/{uid}/allBoard", boardController.GetByUserID)
	myRouter.Get("/sys/allBoard", boardController.GetAllBoard)
	myRouter.Get("/sys/filter", boardController.FilterForSystem)

	// Routers for USER table
	return myRouter
}
