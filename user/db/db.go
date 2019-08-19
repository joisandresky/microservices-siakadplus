package db

import "gopkg.in/mgo.v2"

var mgoSession *mgo.Session

var mongoConnStr = "mongodb://jois:jois123@ds255577.mlab.com:55577/todo_micro_jois"

// var mongoConnStr = "mongodb://localhost:27017"

// GetMongoSession - connection
func GetMongoSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mongoConnStr)
		if err != nil {
			return nil, err
		}
	}
	return mgoSession.Clone(), nil
}
