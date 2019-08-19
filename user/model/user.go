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
	Role     string        `json:"role" bson:"role"`
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

// ShowUser - get one user
func ShowUser(id bson.ObjectId) (User, error) {
	var (
		err error
		u   User
	)
	session, err := db.GetMongoSession()
	if err != nil {
		return u, err
	}
	defer session.Close()

	c := session.DB("todo_micro_jois").C("users")
	err = c.FindId(id).One(&u)

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

	return err
}

// RemoveUser - removing one user
func RemoveUser(id string) error {
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB("todo_micro_jois").C("users")
	err = c.RemoveId(bson.ObjectIdHex(id))

	return err
}
