package main

import (
	"log"
	"net/http"
	"todo/infrastructure"
	"todo/router"
)

func main() {
	log.Println("Database name: ", infrastructure.GetDBName())

	log.Fatal(http.ListenAndServe(":"+infrastructure.GetAppPort(), router.Router()))
}
