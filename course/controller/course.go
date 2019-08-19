package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joisandresky/microservices-siakadplus/course/model"
)

// GetCourseService - get course service handler
func GetCourseService(c *gin.Context) {
	courses, err := model.GetCourses()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An Error Occured to get All Courses",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}
