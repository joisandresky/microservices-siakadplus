package model

import (
	"github.com/joisandresky/microservices-siakadplus/course-old/db"
	"gopkg.in/mgo.v2/bson"
)

// Course - model
type Course struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Semester int           `json:"semester" bson:"semester"`
}

// GetCourses - get all Courses
func GetCourses() (course []Course, err error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return course, err
	}
	defer session.Close()

	c := session.DB("todo_micro_jois").C("courses")
	err = c.Find(nil).All(&course)
	return course, err
}
