package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	majorpb "github.com/joisandresky/microservices-siakadplus/major/proto"
	"google.golang.org/grpc"
)

func main() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = ":8181"
	}
	conn, err := grpc.Dial("msp-major-service:9191", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	majorClient := majorpb.NewMajorServiceClient(conn)
	r := gin.Default()
	r.GET("/api/majors", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &majorpb.ListMajorReq{}
		if resp, err := majorClient.ListMajor(ctx, req); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An Error Occured to get All Majors",
				"err":     err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"majors": resp.Major,
			})
		}
	})

	r.Run(port)
}
