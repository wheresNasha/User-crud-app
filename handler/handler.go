package handler

// handler/user.go : handles user-related HTTP requests.
import (
	"go-crud/initializers"
	"go-crud/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 	Think of Context as a tray passed to the handler function that contains two things:
// The Input (Request): Headers, cookies, URL parameters, and JSON body sent by the user.
// The Output (Response): Tools to write back data to the user.

// Your function uses c.JSON() to write data into that context, which Gin then flushes out to the user as a final HTTP response.

// SELECT * FROM users;
func Retrieve(c *gin.Context) {
	// get all users from the Database : list of all users -> Slice

	// Slice of users go variable
	var users []models.User

	// https://gorm.io/docs/query.html#Retrieving-all-objects
	// Retrive all users: SELECT * FROM users;
	// fills users slice with the data from DB and returns status in result
	result := initializers.DB.Find(&users)

	// GORM Metadata
	// result.RowsAffected // returns found records count, equals `len(users)`
	// result.Error        // returns error

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// return them as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully!",
		"users":   users, // Slice (go variable) is serialized to JSON and displayed to the client
	})
}

// SELECT * FROM users WHERE id = 1;
func RetrieveOne(c *gin.Context) {
	// read the id from the Request URL : Request Param

	// c.Bind(&id) -> not this way as id is in the Request URL and not in the Request Body

	// request URL : http://localhost:3001/retrieve/1
	idStr := c.Param("id") // returns the value of the id parameter
	// int("123") is not valid, so we use strconv.Atoi
	// we can't convert string to int directly using type conversion, so we use strconv.Atoi
	id, err := strconv.Atoi(idStr) // converts the string to an integer
	// https://pkg.go.dev/strconv
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	// get one user from the Database
	var user models.User

	// https://gorm.io/docs/query.html
	// SELECT * FROM users WHERE id = 1;
	result := initializers.DB.First(&user, id)
	// writes the data to the user varibale from the Database
	// first() is a GORM method that returns the first record found or an error

	// result.RowsAffected // returns found records count, equals `len(users)`
	// result.Error        // returns error

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// return them as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully!",
		"users":   user,
	})
}

// HEAD is similar to GET, but without the response body.
func Head(c *gin.Context) {
	// read the id from the Request URL : Request Param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	// SELECT * FROM users WHERE id = 1; -> dont save the result from the query
	result := initializers.DB.First(&models.User{}, id)
	// we are not saving the user, see first param is just the struct {}
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User with id " + strconv.Itoa(id) + " exists",
	})
}

// What HTTP methods are allowed for this resource?
// Browser CORS preflight
// router.OPTIONS("/options", func(c *gin.Context) { // when someone hits /, execute the handler function
// 	c.Header("Allow", "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS")
// 	//Without Access-Control-Allow-Origin, the browser blocks the request — even if status is 200.
// 	c.Status(http.StatusNoContent)
// })

// gin.H{"message": "Hello!"} is one line map no trailing comma
// gin.H{ is a multi-line map with a trailing comma: its a go syntax rule
// 	"message": "Hello!",
// }

// GET /users/search?name=John&min_age=25 : SELECT * FROM users WHERE name LIKE "%John%" AND age >= 25;
// GET /users/search?min_age=25 : SELECT * FROM users WHERE age >= 25;
// GET /users/search?name=John : SELECT * FROM users WHERE name LIKE "%John%";
// GET /users/search : SELECT * FROM users;
// GET /users/search?name=” : SELECT * FROM users WHERE name LIKE "%" (empty string)
// GET /users/search?min_age= : SELECT * FROM users
// both the query parameters are optional
func FindUsersByNameAndAge(c *gin.Context) {

	// get the query parameters from the request URL
	nameStr := c.Query("name")
	ageStr := c.Query("min_age")
	min_age, _ := strconv.Atoi(ageStr)

	// get all users from the Database which match the query parameters : list of all users -> Slice
	var users []models.User

	// first see if name exist or it is just filter by age
	// or it is both name and age or non of them
	// LIKE : https://gorm.io/docs/query.html
	var result *gorm.DB
	if nameStr != "" && ageStr != "" {
		result = initializers.DB.Where("name LIKE ? AND age >= ?", "%"+nameStr+"%", min_age).Find(&users)
	} else if nameStr != "" {
		result = initializers.DB.Where("name LIKE ?", "%"+nameStr+"%").Find(&users)
	} else if ageStr != "" {
		result = initializers.DB.Where("age >= ?", min_age).Find(&users)
	} else {
		result = initializers.DB.Find(&users) // find all users
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// return them as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message":          "Users retrieved successfully!",
		"users":            users,
		"noOfFetchedUsers": result.RowsAffected,
	})
}

func FindUsersByMandatoryNameAndAge(c *gin.Context) {
	var users []models.User
	var query struct {
		Name   string `form:"name" binding:"required"`
		MinAge int    `form:"min_age" binding:"required"`
	}
	err := c.ShouldBindQuery(&query)
	// It works like ShouldBindJSON but binds query parameters from the URL.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := initializers.DB.
		Where("name LIKE ? AND age >= ?", "%"+query.Name+"%", query.MinAge).
		Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// return them as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully!",
		"users":   users,
	})
}

// INSERT INTO users (name, email) VALUES ("Sakshi", "sakshi@gmail.com");
func Create(c *gin.Context) {
	// get data from request body
	// adding tags for Input validation on Create
	// required : must be present
	// email : must be a valid email address
	// min=0 : must be greater than or equal to 0
	var body struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Age   int    `json:"age" binding:"min=0"`
	}

	err := c.ShouldBindJSON(&body) // returns 400 if missing
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// read the name and email from the request body
	// POST with {"name":"Sakshi", "email":"sakshi@gmail.com"}
	//c.Bind(&body) // https://pkg.go.dev/github.com/gin-gonic/gin@v1.12.0#Context.Bind
	// binds the request body from client to the body struct in the handler

	// create a new user in the Database
	// https://gorm.io/docs/create.html
	// 	user := models.User{Name: "Sakshi", Email: "sakshi@gmail.com"}
	user := models.User{Name: body.Name, Email: body.Email, Age: body.Age}

	// Create a single record
	result := initializers.DB.Create(&user) // pass pointer of data to Create

	// result has two fields : Error and RowsAffected
	// Error is the error message if the operation fails
	// RowsAffected is the number of rows affected by the operation
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// return the created user
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully!",
		"user":    user,
	})
}

func Update(c *gin.Context) {

	// get which id to update from the request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	// Update details from the request body

	var body struct {
		Name  string
		Email string
		Age   int
	}
	c.Bind(&body)
	// binds the request body from client to the body struct in the handler

	// https://gorm.io/docs/update.html

	var user models.User
	initializers.DB.First(&user, id) // get the user from the DB
	// SELECT * FROM users WHERE id = 1;

	user.Name = body.Name
	user.Email = body.Email
	user.Age = body.Age
	initializers.DB.Save(&user) // now update the user in the DB

	// Save is an upsert function: it updates the record if it exists, otherwise it creates a new record
	// UPDATE users SET name='jinzhu 2', email='dijw@gmail.com' age=100', updated_at = '2013-11-17 21:34:10'
	// WHERE id=111;

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated user details successfully for id :" + strconv.Itoa(id),
		"user":    user,
	})
}

// update a particular field(age) of a resource by id
func Patch(c *gin.Context) {
	// get which id to update from the request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	// Update details from the request body

	var body struct {
		Age int
	}
	c.Bind(&body)
	// binds the request body from client to the body struct in the handler(age only)

	// https://gorm.io/docs/update.html

	var user models.User
	initializers.DB.First(&user, id) // get the user from the DB

	user.Age = body.Age
	initializers.DB.Save(&user) // now update the user in the DB
	c.JSON(http.StatusOK, gin.H{
		"message": "Updated age column of user with id :" + strconv.Itoa(id) + " successfully",
		"user":    user,
	})
}

func Delete(c *gin.Context) {

	// get which id to delete from the request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	// Delete the user from the Database
	// https://gorm.io/docs/delete.html

	result := initializers.DB.Delete(&models.User{}, id)
	// DELETE FROM users WHERE id = 10;

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully for id :" + strconv.Itoa(id),
	})
}
