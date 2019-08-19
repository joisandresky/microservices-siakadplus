package controller

import (
	"log"
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

// ShowUserService - show one user service handler
func ShowUserService(c *gin.Context) {
	paramID := c.Param("id")
	user, err := model.ShowUser(paramID)
	log.Println(err.Error())
	if err.Error() == "not found" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "User not found with ID:" + paramID,
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An Error Occured to show user ID:" + paramID,
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// AddUserService - add one user service handler
func AddUserService(c *gin.Context) {
	newUser := model.User{
		ID:       bson.NewObjectId(),
		Name:     c.PostForm("name"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
		Role:     c.PostForm("role"),
	}

	err := model.AddUser(newUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An Error Occured to Create User",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created!",
		"user":    newUser,
	})
}

// RemoveUserService - remove one user service handler
func RemoveUserService(c *gin.Context) {
	param := c.Param("id")

	err := model.RemoveUser(param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An Error Occured to Remove User",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Deleted!",
		"id":      param,
	})
}
