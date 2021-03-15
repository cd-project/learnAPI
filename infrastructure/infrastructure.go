package infrastructure

import (
	"log"
	"todo/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db         *gorm.DB
	appPort    string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
)

func loadParameters() {
	appPort = "8080"
	dbHost = "localhost"
	dbPort = "5432"
	dbUser = "postgres"
	dbPassword = "zxc1234567890A"
	//dbPassword = "0000" // My password
	dbName = "model"
}

func openConnection() (*gorm.DB, error) {
	connectSQL := "host=" + dbHost +
		" user=" + dbUser +
		" dbname=" + dbName +
		" password=" + dbPassword +
		" sslmode=disable"

	db, err := gorm.Open("postgres", connectSQL)
	if err != nil {
		log.Printf("Has problem at connection database: %+v\n", err)
		return nil, err
	}

	return db, nil
}

func closeConnection(db *gorm.DB) {
	_ = db.Close()
}

// InitDatabase open connection and migrate database
func InitDatabase() error {
	// connect db
	var err error
	db, err = openConnection()
	if err != nil {
		return err
	}
	// Migrating db
	log.Println("Migrating database...")
	db.AutoMigrate(
		&model.Todo{},
	)

	return nil
}

func init() {
	// load param
	loadParameters()
	// init database
	InitDatabase()
}

// GetDB export db
func GetDB() *gorm.DB {
	return db
}

// Get db name
func GetDBName() string {
	return dbName
}

// GetAppPort export app port
func GetAppPort() string {
	return appPort
}
