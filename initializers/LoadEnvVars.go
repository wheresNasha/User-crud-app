package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVars() { // function name is capitalized to make it public
	// If you call Load without any args it will default to loading .env in the current path.
	// .env file is loaded into the environment variables
	err := godotenv.Load()
	// if .env file is not found, it will log a message and continue execution
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
