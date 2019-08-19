package model

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/joisandresky/microservices-siakadplus/user/db"
)

// User - user model
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Email    string        `json:"email" bson:"email"`
	Password string        `json:"password" bson:"password"`
}

// GetUsers - get all users
func GetUsers() (u []User, err error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return u, err
	}
	defer session.Close()

	c := session.DB("todo_micro_jois").C("users")
	err = c.Find(nil).All(&u)
	return u, err
}

// AddUser - add one user
func AddUser(u User) error {
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}

	defer session.Close()

	c := session.DB("todo_micro_jois").C("users")
	err = c.Insert(&u)

	return nil
}
