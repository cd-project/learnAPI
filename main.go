package main

import (
	"log"
	"net/http"
	"todo/infrastructure"
	"todo/router"
)

// @title Swagger
// @version 1.0
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Println("Database name: ", infrastructure.GetDBName())

	log.Fatal(http.ListenAndServe(":"+infrastructure.GetAppPort(), router.Router()))

}
