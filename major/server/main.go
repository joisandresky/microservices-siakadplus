package main

import (
	"context"
	"log"
	"net"

	"github.com/joisandresky/microservices-siakadplus/major/db"
	majorpb "github.com/joisandresky/microservices-siakadplus/major/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

type server struct{}

func main() {
	// var port = os.Getenv("PORT")
	// if port == "" {
	// 	port = ":8181"
	// }
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
	majors := []*majorpb.Major{}
	session, err := db.GetMongoSession()
	if err != nil {
		return &majorpb.ListMajorRes{}, err
	}
	defer session.Close()

	c := session.DB("todo_micro_jois").C("majors")
	err = c.Find(bson.M{}).All(&majors)

	return &majorpb.ListMajorRes{Major: majors}, err
}

func (s *server) CreateMajor(ctx context.Context, req *majorpb.CreateMajorReq) (*majorpb.CreateMajorRes, error) {
	return &majorpb.CreateMajorRes{Major: &majorpb.Major{}}, nil
}

func (s *server) ReadMajor(ctx context.Context, req *majorpb.ReadMajorReq) (*majorpb.ReadMajorRes, error) {
	_ = req.GetId()

	return &majorpb.ReadMajorRes{Major: &majorpb.Major{}}, nil
}

func (s *server) UpdateMajor(ctx context.Context, req *majorpb.UpdateMajorReq) (*majorpb.UpdateMajorRes, error) {
	return &majorpb.UpdateMajorRes{Major: &majorpb.Major{}}, nil
}

func (s *server) DeleteMajor(ctx context.Context, req *majorpb.DeleteMajorReq) (*majorpb.DeleteMajorRes, error) {
	_ = req.GetId()

	return &majorpb.DeleteMajorRes{Success: true}, nil
}
