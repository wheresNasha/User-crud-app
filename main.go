package main

// main.go : the entry point of the application.
// start server, wire things together

import (
	"go-crud/initializers"
	"go-crud/migrator"
	"go-crud/routes"

	"github.com/gin-gonic/gin"
)

func init() { // init is a special function that is called before main is executed
	initializers.LoadEnvVars()
	initializers.ConnectToDB()

	// AutoMigrate will create tables (users), missing foreign keys, and create missing indexes
	migrator.Migrate()
}

func main() {
	router := gin.Default() // returns a Engine instance which contains the router

	// runs on every request before routes are executed. It is a middleware.
	// router.Use(middleware.CORS())

	// wire all routes to their handlers
	// we can directly do
	// router.GET("/retrieve", handler.Retrieve)
	// but it gets messy fast : and at scale (50 routes), it becomes difficult to maintain
	// so we use routes.go to wire all routes to their handlers
	// routes.go is a better way to organize routes

	routes.Register(router)

	router.Run() // listens on 0.0.0.0:8080 by default and serves the application
	// Run internally uses ListenAndServe() to start the server on the port specified in the router.
}
