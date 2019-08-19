package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joisandresky/microservices-siakadplus/user/controller"
)

func main() {
	r := gin.Default()

	r.GET("/api/users", controller.GetUserService)
	r.POST("/api/users", controller.AddUserService)

	var port = os.Getenv("PORT")
	if port == "" {
		port = ":8181"
	}

	r.Run(port)
}
