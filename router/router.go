package router

import (
	"log"
	"net/http"
	"os"
	"todo/controller"
	"todo/infrastructure"

	_ "todo/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
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
	userController := controller.NewUserController()

	// TEST
	myRouter.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	// Middleware routing
	myRouter.Route("/", func(router chi.Router) {
		// public routes
		myRouter.Post("/user/login", userController.Login)
		myRouter.Post("/user/login/token", userController.LoginWithToken)
		router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		// protected route
		router.Group(func(protectedRoute chi.Router) {
			// Middleware authentication
			protectedRoute.Use(jwtauth.Verifier(infrastructure.GetEncodeAuth()))
			protectedRoute.Use(jwtauth.Authenticator)

			// Todo table
			protectedRoute.Route("/work", func(subRouter chi.Router) {
				subRouter.Post("/create", todoController.Create)
				subRouter.Get("/search/{id}", todoController.GetByID)
				subRouter.Put("/updater/{id}", todoController.Update)
				subRouter.Delete("/delete/{id}", todoController.Delete)
			})
		})

	})

	// protected route

	// Routers for TODOS table
	myRouter.Post("/work/create", todoController.Create)
	myRouter.Get("/work/all", todoController.GetAll)
	myRouter.Get("/work/search/{id}", todoController.GetByID)
	myRouter.Put("/work/updater/{id}", todoController.Update)
	myRouter.Delete("/work/delete/{id}", todoController.Delete)

	// Routers for BOARD table
	myRouter.Post("/board/{uid}/create", boardController.CreateBoard)
	myRouter.Put("/board/{boardid}/update", boardController.UpdateBoard)
	myRouter.Delete("/board/delete/{boardid}", boardController.DeleteBoard)
	myRouter.Get("/board/{uid}/allBoard", boardController.GetByUserID)
	myRouter.Get("/sys/allBoard", boardController.GetAllBoard)
	myRouter.Put("/sys/filter", boardController.Filter)

	// Routers for USER table
	myRouter.Get("/user/all", userController.GetAll)
	myRouter.Post("/user/create", userController.CreateUser)
	myRouter.Get("/user/{uid}", userController.GetByID)
	myRouter.Put("/user/{uid}/modify/pwd", userController.ChangePassword)
	myRouter.Put("/user/{uid}/modify/role", userController.ChangeRole)
	myRouter.Post("/user/login", userController.Login)
	myRouter.Post("/user/login/token", userController.LoginWithToken)
	return myRouter
}
