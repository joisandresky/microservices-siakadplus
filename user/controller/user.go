package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joisandresky/microservices-siakadplus/user/model"
	"gopkg.in/mgo.v2/bson"
)

// GetUserService - get user service handler
func GetUserService(c *gin.Context) {
	users, err := model.GetUsers()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An Error Occured to get All Users",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// AddUserService - add one user service handler
func AddUserService(c *gin.Context) {
	newUser := model.User{
		ID:       bson.NewObjectId(),
		Name:     c.PostForm("name"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	err := model.AddUser(newUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An Error Occured to get All Users",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created!",
		"user":    newUser,
	})
}
