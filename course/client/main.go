package main

import (
	"context"
	"github.com/gin-gonic/gin"
	coursepb "github.com/joisandresky/microservices-siakadplus/course/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = ":8181"
	}

	gConn := os.Getenv("g_course_server")
	if gConn == "" {
		gConn = "localhost"
	}
	conn, err := grpc.Dial(gConn+":9292", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	courseClient := coursepb.NewCourseServiceClient(conn)
	r := gin.Default()
	r.GET("/api/courses", CORSMiddleware(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			handleError(c, err, "An error to parsing Page/Limit", nil)
		}

		req := &coursepb.ListCourseReq{
			Page:  int32(page),
			Limit: int32(limit),
		}
		if resp, err := courseClient.ListCourse(ctx, req); err != nil {
			handleError(c, err, "An Error Occured to get All Courses", nil)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"courses": resp.Course,
				"total":   resp.Total,
			})
		}
	})

	r.GET("/api/courses/:id", CORSMiddleware(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()
		req := &coursepb.ReadCourseReq{
			Id: c.Param("id"),
		}
		resp, err := courseClient.ReadCourse(ctx, req)
		if err != nil {
			handleError(c, err, "An Error Occured To Read A Course", nil)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"course": resp.Course,
			})
		}

	})

	r.POST("/api/courses", CORSMiddleware(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &coursepb.CreateCourseReq{
			Course: &coursepb.Course{
				Name:     c.PostForm("name"),
				Semester: c.PostForm("semester"),
			},
		}
		if resp, err := courseClient.CreateCourse(ctx, req); err != nil {
			handleError(c, err, "An Error Occured to create a Course", nil)
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Course Created!",
				"course":  resp.Course,
			})
		}
	})

	r.PUT("/api/courses/:id", CORSMiddleware(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &coursepb.UpdateCourseReq{
			Course: &coursepb.Course{
				Id:       c.Param("id"),
				Name:     c.PostForm("name"),
				Semester: c.PostForm("semester"),
			},
		}

		if resp, err := courseClient.UpdateCourse(ctx, req); err != nil {
			handleError(c, err, "An Error Occured when Updating a Course", nil)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Course Updated!",
				"course":  resp.Course,
			})
		}
	})

	r.DELETE("/api/courses/:id", CORSMiddleware(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		req := &coursepb.DeleteCourseReq{
			Id: c.Param("id"),
		}
		if resp, err := courseClient.DeleteCourse(ctx, req); err != nil {
			handleError(c, err, "An Error Occured when Delete a Course", resp.GetSuccess())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Course Deleted!",
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
			"data":    data,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": message,
			"err":     err,
			"data":    data,
		})
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
