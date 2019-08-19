package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joisandresky/microservices-siakadplus/user/controller"
)

func main() {
	r := gin.Default()

	r.GET("/api/users", controller.GetUserService)
	r.GET("/api/users/:id", controller.ShowUserService)
	r.POST("/api/users", controller.AddUserService)
	r.DELETE("/api/users/:id", controller.RemoveUserService)

	var port = os.Getenv("PORT")
	if port == "" {
		port = ":8181"
	}

	r.Run(port)
}
