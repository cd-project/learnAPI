package router

import (
	"net/http"
	"todo/controller"
	"todo/infrastructure"

	_ "todo/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
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
		myRouter.Post("/user/create", userController.CreateUser)
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
				subRouter.Get("/all", todoController.GetAll)
			})

			// Board table
			protectedRoute.Route("/board", func(subRouter chi.Router) {
				subRouter.Post("/{uid}/create", boardController.CreateBoard)
				subRouter.Put("/{boardid}/update", boardController.UpdateBoard)
				subRouter.Delete("/delete/{boardid}", boardController.DeleteBoard)
				subRouter.Get("/{uid}/allBoard", boardController.GetByUserID)
				subRouter.Get("/allBoard", boardController.GetAllBoard)
				subRouter.Put("/filter", boardController.Filter)
			})

			// User table
			protectedRoute.Route("/user", func(subRouter chi.Router) {
				subRouter.Get("/all", userController.GetAll)
				subRouter.Get("/{uid}", userController.GetByID)
				subRouter.Put("/{uid}/modify/pwd", userController.ChangePassword)
				subRouter.Put("/{uid}/modify/role", userController.ChangeRole)
			})
		})

	})

	return myRouter
}
