package main

import "C"
import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"

	"github.com/joisandresky/microservices-siakadplus/major/db"
	majorpb "github.com/joisandresky/microservices-siakadplus/major/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

type Major struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name 	string             `bson:"name"`
	Level   string             `bson:"level"`
	Head    string             `bson:"head"`
	Status  string				`bson:"status"`
}

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":9191")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	majorpb.RegisterMajorServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	} else {
		log.Printf("gRPC Server Running [MajorService]")
	}
}

func (s *server) ListMajor(ctx context.Context, req *majorpb.ListMajorReq) (*majorpb.ListMajorRes, error) {
	var majors []*majorpb.Major
	data  := &Major{}
	session, err := db.GetMongoSession()
	if err != nil {
		return &majorpb.ListMajorRes{}, err
	}

	findOptions := options.Find()
	collection := session.Database("todo_micro_jois").Collection("majors")

	curr, err := collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		return &majorpb.ListMajorRes{}, err
	}

	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		err := curr.Decode(data)
		if err != nil {
			log.Fatalln(err)
		}
		majors = append(majors, &majorpb.Major{
			Id:                   data.ID.Hex(),
			Name:                 data.Name,
			Level:                data.Level,
			Head:                 data.Head,
			Status:               data.Status,
		})
	}

	return &majorpb.ListMajorRes{Major: majors}, err
}

func (s *server) CreateMajor(ctx context.Context, req *majorpb.CreateMajorReq) (*majorpb.CreateMajorRes, error) {
	major := &majorpb.Major{
		//Id: 				pmongo.NewObjectId(primitive.NewObjectID()),
		Name:                 req.GetMajor().Name,
		Level:                req.GetMajor().Level,
		Head:                 req.GetMajor().Head,
		Status:               req.GetMajor().Status,
	}
	session, err := db.GetMongoSession()
	if err != nil {
		return &majorpb.CreateMajorRes{}, err
	}
	collection := session.Database("todo_micro_jois").Collection("majors")

	result, err := collection.InsertOne(context.TODO(), major)
	log.Println("result insert", result)

	oid := result.InsertedID.(primitive.ObjectID)
	major.Id = oid.Hex()

	return &majorpb.CreateMajorRes{Major: major}, err
}

func (s *server) ReadMajor(ctx context.Context, req *majorpb.ReadMajorReq) (*majorpb.ReadMajorRes, error) {
	majorID := req.GetId()
	id, err := primitive.ObjectIDFromHex(majorID)
	if err != nil {
		return &majorpb.ReadMajorRes{Major: &majorpb.Major{}}, err
	}

	major  := &Major{}
	session, err := db.GetMongoSession()
	if err != nil {
		return &majorpb.ReadMajorRes{Major: &majorpb.Major{}}, err
	}
	collection := session.Database("todo_micro_jois").Collection("majors")
	query := bson.M{"_id": id}

	err = collection.FindOne(context.TODO(), query).Decode(major)

	return &majorpb.ReadMajorRes{Major: &majorpb.Major{
		Id:                   major.ID.Hex(),
		Name:                 major.Name,
		Level:                major.Level,
		Head:                 major.Head,
		Status:               major.Status,
	}}, err
}

func (s *server) UpdateMajor(ctx context.Context, req *majorpb.UpdateMajorReq) (*majorpb.UpdateMajorRes, error) {
	major := req.GetMajor()

	id, err := primitive.ObjectIDFromHex(major.GetId())
	if err != nil {
		return &majorpb.UpdateMajorRes{Major: &majorpb.Major{}}, err
	}

	update := bson.M{
		"name": major.GetName(),
		"level": major.GetLevel(),
		"head": major.GetHead(),
		"status": major.GetStatus(),
	}

	filter := bson.M{ "_id": id}
	session, err := db.GetMongoSession()
	if err != nil {
		return &majorpb.UpdateMajorRes{Major: &majorpb.Major{}}, err
	}
	collection := session.Database("todo_micro_jois").Collection("majors")

	result := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
	decoded := Major{}
	err = result.Decode(&decoded)
	if err != nil {
		return &majorpb.UpdateMajorRes{Major: &majorpb.Major{}}, err
	}

	return &majorpb.UpdateMajorRes{Major: &majorpb.Major{
		Id: decoded.ID.Hex(),
		Name: decoded.Name,
		Level: decoded.Level,
		Head: decoded.Head,
		Status: decoded.Status,
	}}, nil
}

func (s *server) DeleteMajor(ctx context.Context, req *majorpb.DeleteMajorReq) (*majorpb.DeleteMajorRes, error) {
	majorID := req.GetId()
	id, err  := primitive.ObjectIDFromHex(majorID)
	if err != nil {
		return &majorpb.DeleteMajorRes{Success:false}, err
	}
	session, err := db.GetMongoSession()
	if err != nil {
		return &majorpb.DeleteMajorRes{Success:false}, err
	}
	collection := session.Database("todo_micro_jois").Collection("majors")

	result, err := collection.DeleteOne(context.TODO(), bson.M{ "_id": id })
	log.Println("result delete", result)

	return &majorpb.DeleteMajorRes{Success: true}, err
}
