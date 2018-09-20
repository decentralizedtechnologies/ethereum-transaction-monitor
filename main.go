package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	database "gitlab.com/canya-com/canwork-database-client"
	GET "gitlab.com/canya-com/canya-ethereum-tx-api/internal/GET"
	POST "gitlab.com/canya-com/canya-ethereum-tx-api/internal/POST"
	"google.golang.org/appengine"
)

var (
	// DatabaseClient : global gorm.DB instance
	DatabaseClient *gorm.DB
	// DefaultTimeout : 1 day
	DefaultTimeout = 1
)

func init() {
	loadEnvironmentFile()
	makeDatabaseConnection()

	http.HandleFunc("/tx/details", GET.TransactionDetails)
	http.HandleFunc("/tx/monitor", POST.MonitorTransaction)
}

func main() {
	defer DatabaseClient.Close()
	appengine.Main()
}

func loadEnvironmentFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func makeDatabaseConnection() {
	var err error
	dsn := makeDsnString()
	DatabaseClient, err = database.NewDatabaseClient("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func makeDsnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		GetEnv("DB_CANWORK_STAGING_USERNAME", ""),
		GetEnv("DB_CANWORK_STAGING_PWD", ""),
		GetEnv("DB_CANWORK_STAGING_HOST", ""),
		GetEnv("DB_CANWORK_STAGING_PORT", ""),
		GetEnv("DB_CANWORK_STAGING_NAME", ""))
}
