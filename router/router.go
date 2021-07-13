package router

import (
	"log"
	"net/http"
	"todo/controller"
	"todo/infrastructure"
	"todo/middlewares"

	_ "todo/docs"

	"github.com/casbin/casbin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
)

func CasbinEnforce(e *casbin.Enforcer, role, url, method string) {
	chek := e.Enforce(role, url, method)
	if chek {
		log.Printf("Allow %s %s %s\n", role, url, method)
	} else {
		log.Printf("Deny %s %s %s\n", role, url, method)
	}

}
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
			// e := casbin.NewEnforcer("./infrastructure/authz_model.conf", "./infrastructure/authz_policy.csv")
			protectedRoute.Use(middlewares.Authorizer(infrastructure.GetEnforce()))

			// Todo table < a little bit deprecated >
			protectedRoute.Route("/work", func(subRouter chi.Router) {
				subRouter.Post("/create", todoController.Create)        // user
				subRouter.Get("/search/{id}", todoController.GetByID)   // user
				subRouter.Put("/update/{id}", todoController.Update)    // user
				subRouter.Delete("/delete/{id}", todoController.Delete) // user
				subRouter.Get("/all", todoController.GetAll)            // admin
			})

			// Board table
			protectedRoute.Route("/board", func(subRouter chi.Router) {
				subRouter.Post("/{uid}/create", boardController.CreateBoard)       // user
				subRouter.Put("/{boardid}/update", boardController.UpdateBoard)    // user
				subRouter.Delete("/delete/{boardid}", boardController.DeleteBoard) // user
				subRouter.Get("/{uid}/allBoard", boardController.GetByUserID)      // user
				subRouter.Get("/all", boardController.GetAllBoard)                 // admin
				subRouter.Put("/filter", boardController.Filter)                   // user
			})

			// User table
			protectedRoute.Route("/user", func(subRouter chi.Router) {
				subRouter.Get("/all", userController.GetAll)                 // admin
				subRouter.Get("/{uid}", userController.GetByID)              // admin
				subRouter.Put("/modify/pwd", userController.ChangePassword)  // user
				subRouter.Put("/modify/role", userController.ChangeRole)     // user
				subRouter.Put("/reset/{uid}", userController.ResetPassword)  // user
				subRouter.Delete("/delete/{uid}", userController.DeleteUser) // admin
			})
		})

	})

	return myRouter
}
