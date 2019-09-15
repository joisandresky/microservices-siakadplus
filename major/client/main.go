package main

import (
	"context"
	"google.golang.org/grpc/status"
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

	gConn := os.Getenv("g_major_server")
	if gConn == "" {
		gConn = "localhost"
	}
	conn, err := grpc.Dial(gConn + ":9191", grpc.WithInsecure())
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
			handleError(c, err, "An Error Occured to get All Majors", nil)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"majors": resp.Major,
			})
		}
	})

	r.GET("/api/majors/:id", func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()
		req := &majorpb.ReadMajorReq{
			Id:                   c.Param("id"),
		}
		resp, err := majorClient.ReadMajor(ctx, req);
		if  err != nil {
			handleError(c, err, "An Error Occured To Read A Major", nil)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"major": resp.Major,
			})
		}

	})

	r.POST("/api/majors", func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &majorpb.CreateMajorReq{
			Major: &majorpb.Major{
				Name:   c.PostForm("name"),
				Level:  c.PostForm("level"),
				Head:   c.PostForm("head"),
				Status: c.PostForm("status"),
			},
		}
		if resp, err := majorClient.CreateMajor(ctx, req); err != nil {
			handleError(c, err, "An Error Occured to create a Major", nil)
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Major Created!",
				"major": resp.Major,
			})
		}
	})

	r.PUT("/api/majors/:id", func(c * gin.Context){
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &majorpb.UpdateMajorReq{
			Major: &majorpb.Major{
				Id: c.Param("id"),
				Name:   c.PostForm("name"),
				Level:  c.PostForm("level"),
				Head:   c.PostForm("head"),
				Status: c.PostForm("status"),
			},
		}

		if resp, err := majorClient.UpdateMajor(ctx, req); err != nil {
			handleError(c, err, "An Error Occured when Updating a Major", nil)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Major Updated!",
				"major": resp.Major,
			})
		}
	})

	r.DELETE("/api/majors/:id", func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &majorpb.DeleteMajorReq{
			Id: c.Param("id"),
		}
		if resp, err := majorClient.DeleteMajor(ctx, req); err != nil {
			handleError(c, err, "An Error Occured when Delete a Major", resp.GetSuccess())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Major Deleted!",
				"success": resp.Success,
			})
		}
	})

	r.Run(port)
}

func handleError(c *gin.Context, err error, message string, data interface{}) {
	ok := status.Convert(err)

	if ok.Message() == "mongo: no documents in result" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found!",
			"err":     err,
			"data": data,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": message,
			"err": err,
			"data": data,
		})
	}
}
