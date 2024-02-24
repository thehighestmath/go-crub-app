package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/gin-swagger"
	_ "web/docs"
	swaggerfiles "github.com/swaggo/files"
)

type User struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

var Users = []User{
	{
		ID:   1,
		Name: "Alice",
	},
	{
		ID:   2,
		Name: "Bob",
	},
}

// @Summary Get all users
// @Description Get all users from the database
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Router /users [get]
func ViewUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Users)
}

func ViewUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "invalid user ID",
		})
	}
	for _, user := range Users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "fail",
		"message": "not found",
	})
}

func AddUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid JSON data",
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	Users = append(Users, u)
	log.Println("Пользователь успешно добавлен:", u)

	c.JSON(http.StatusOK, gin.H{"message": "user added successfully"})
}

func UpdateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid JSON data",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	for i, user := range Users {
		if user.ID == u.ID {
			Users[i] = u
			c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "fail", "message": "user not found"})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "invalid user ID",
		})
	}

	index := -1
	for i, user := range Users {
		if user.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "user not found",
		})
		return
	}

	Users = append(Users[:index], Users[index + 1:]...)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/users", ViewUsers)
	r.GET("/users/:id", ViewUser)
	r.POST("/users/add", AddUser)
	r.PUT("/users/update", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
