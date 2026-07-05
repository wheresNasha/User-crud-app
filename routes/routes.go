package routes

// routes.go : wires all HTTP routes to their handlers.
//  WHAT URLs exist?
//  WHAT HTTP methods are allowed for each URL?
//  WHAT handler functions are called for each URL?
//  WHAT data is sent to and received from each URL?
//  WHAT HTTP status codes are returned for each URL?
//  WHAT HTTP headers are sent to and received from each URL?
//  WHAT HTTP body is sent to and received from each URL?
//  WHAT HTTP cookies are sent to and received from each URL?
//  WHAT HTTP authentication is required for each URL?

import (
	"go-crud/handler"

	"github.com/gin-gonic/gin"
)

// Register wires all HTTP routes to their handlers.
func Register(router *gin.Engine) {
	// you register the handlers here
	// but never call them right away like its not handler.Retrieve() here
	// When user hits GET /retrieve  →  Gin calls handler.Retrieve(c)
	// You register handlers at startup; Gin invokes them per request.

	// when someone hits /retrieve, execute the handler function
	// Delete() sets deleted_at, doesn't hard-delete
	// Find() auto-excludes soft-deleted rows
	router.GET("/users", handler.Retrieve)

	//https://gin-gonic.com/en/docs/routing/param-in-path/
	router.GET("/users/:id", handler.RetrieveOne)

	// create a new resource
	router.POST("/users", handler.Create)

	// update the entire resource by id
	router.PUT("/users/:id", handler.Update)

	// update a particular field(age) of a resource by id
	router.PATCH("/users/:id", handler.Patch)

	// delete the entire resource by id
	router.DELETE("/users/:id", handler.Delete)

	// we dont do handler.Patch() here because we want to call it when PATCH /patch arrives
	// handler.Patch    // "call this when PATCH /patch arrives"
	// handler.Patch(c *gin.Context)  // "call this RIGHT NOW"

	//router.GET(path string, handler func(c *gin.Context))
	// It wants a function reference — something to call later when a request hits /retrieve.

	// HEAD is similar to GET, but without the response body.
	// It's useful for checking if a resource exists without downloading the content.
	// Postman shows the status (200) but no message because HTTP forbids a body on HEAD responses. Gin honors that even if you call c.JSON().
	//
	// Think of HEAD as a quick "ping" to check if a file exists without downloading it.
	// Purpose is to read metadata like size, last modified date, etc. without downloading the content.
	// Health checks on large objects without downloading the content.
	// "Does this file exist? How big is it?"
	router.HEAD("/users/:id", handler.Head)

	// https://gin-gonic.com/en/docs/routing/querystring-param/
	// GET /users/search?name=John&min_age=25 -> result has name John and age >= 25
	// GET /users/search?name=John -> result has name John
	// GET /users/search?min_age=25 -> result has age >= 25
	router.GET("/users/search", handler.FindUsersByNameAndAge)

	// GET /users/search?name=John&min_age=25
	router.GET("/users/search2", handler.FindUsersByMandatoryNameAndAge)
}
