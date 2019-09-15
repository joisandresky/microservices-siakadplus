package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var mgoSession *mgo.Session

var mongoConnStr = "mongodb://jois:jois123@ds255577.mlab.com:55577/todo_micro_jois?retryWrites=false"

// var mongoConnStr = "mongodb://localhost:27017"

// GetMongoSession - connection
func GetMongoSession() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoConnStr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}
