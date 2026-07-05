package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB connection singleton
var DB *gorm.DB // global variable to store the database connection
// initializers.DB — one connection pool, reused across requests.

func ConnectToDB() {
	// https://gorm.io/docs/connecting_to_the_database.html
	// I am using neon pg db : https://console.neon.tech/app/projects/proud-pond-04715218
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}
