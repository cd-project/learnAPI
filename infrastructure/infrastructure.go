package infrastructure

import (
	"crypto/rsa"
	"log"
	"todo/model"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	APPPORT    = ""
	DBHOST     = ""
	DBPORT     = ""
	DBUSER     = ""
	DBPASSWORD = ""
	DBNAME     = ""

	HTTPSWAGGER = ""
	ROOTPATH    = ""

	RSAPUBLICPATH      = ""
	RSAPRIVATEPATH     = ""
	RSAPRIVATEPASSWORD = ""

	EXTENDHOUR        = ""
	EXTENDHOURREFRESH = ""

	NANO_TO_SECOND = 1000000000
	Extend_Hour    = 72
)

var (
	db         *gorm.DB
	appPort    string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	encodeAuth *jwtauth.JWTAuth
	decodeAuth *jwtauth.JWTAuth
	privateKey *rsa.PrivateKey
	publicKey  interface{}

	InfoLog *log.Logger
	ErrLog  *log.Logger

	rsaPublicPath  string
	rsaPrivatePath string

	extendHour        int
	extendHourRefresh int
)

func loadParameters() {
	appPort = "8080"
	dbHost = "localhost"
	dbPort = "5432"
	dbUser = "postgres"
	dbPassword = "0000" // My password
	dbName = "learnAPI"
	extendHour = 72
	extendHourRefresh = 1440
	rsaPublicPath = "./infrastructure/public_key.pem"
	rsaPrivatePath = "./infrastructure/private_key.pem"
}

// open connection to database
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
		&model.Board{},
		&model.User{},
		&model.Profile{},
	)

	return nil
}

func init() {
	// load param
	loadParameters()
	// init database
	InitDatabase()
	if err := loadAuthToken(); err != nil {
		ErrLog.Println(err)
	}
}

func GetEnforce() *casbin.Enforcer {
	enforcer, _ := casbin.NewEnforcer("./infrastructure/authz_model.conf",
		"./infrastructure/authz_policy.csv")

	return enforcer
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

// GetEncodeAuth get token auth
func GetEncodeAuth() *jwtauth.JWTAuth {
	return encodeAuth
}

// GetExtendAccessMinute export access extend minute
func GetExtendAccessHour() int {
	return extendHour
}

// GetExtendRefreshHour export refresh extends hour
func GetExtendRefreshHour() int {
	return extendHourRefresh
}

//
func GetPublicKey() interface{} {
	return publicKey
}
