package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joisandresky/microservices-siakadplus/course-old/controller"
)

func main() {
	r := gin.Default()

	r.GET("/api/courses", controller.GetCourseService)

	var port = os.Getenv("PORT")
	if port == "" {
		port = ":8181"
	}

	r.Run(port)
}
