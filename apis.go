package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

type Todo struct {
	gorm.Model
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	Done        bool   `gorm:"column:finished"`
}

func InitialMigration() {
	db, err := gorm.Open("postgres", "user=postgres password=zxc1234567890A dbname=model sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	db.AutoMigrate(&Todo{})
}
func getAgenda(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "user=postgres password=zxc1234567890A dbname=model sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var todos []Todo
	db.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

func newWork(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "newWork endpoint")
	db, err := gorm.Open("postgres", "user=postgres password=zxc1234567890A dbname=model sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)             // request variable
	strID := vars["id"]             // id as a string
	title := vars["title"]          // title request
	intID, _ := strconv.Atoi(strID) // cvt id to int
	db.Create(&Todo{ID: intID, Title: title})

	fmt.Fprintf(w, "New work succesfully added!")
}

func deleteWork(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "user=postgres password=zxc1234567890A dbname=model sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	intID := vars["id"]

	var todo Todo
	db.Where("id = ?", intID).Find(&todo)
	db.Delete(&todo)
	fmt.Fprintf(w, "Deleted")

}

func updateWork(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "user=postgres password=zxc1234567890A dbname=model sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	intID := vars["id"]
	title := vars["title"]

	var todo Todo
	db.Where("id = ?", intID).Find(&todo)

	todo.Title = title

	db.Save(&todo)
	fmt.Fprintf(w, "Updated Successfully!")
}
