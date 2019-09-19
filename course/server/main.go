package main

import (
	"context"
	"github.com/joisandresky/microservices-siakadplus/course/db"
	coursepb "github.com/joisandresky/microservices-siakadplus/course/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Course struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name 	string             `bson:"name"`
	Semester   string             `bson:"semester"`
}

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":9191")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	coursepb.RegisterCourseServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	} else {
		log.Printf("gRPC Server Running [CourseService]")
	}
}

func (s *server) ListCourse(ctx context.Context, req *coursepb.ListCourseReq) (*coursepb.ListCourseRes, error) {
	var courses []*coursepb.Course
	data  := &Course{}
	page := req.GetPage()
	limit := int64(req.GetLimit())
	skip := int64(page - 1) * limit

	session, err := db.GetMongoSession()
	if err != nil {
		return &coursepb.ListCourseRes{}, err
	}

	findOptions := options.Find().SetLimit(limit).SetSkip(skip)
	countOptions := options.Count()
	collection := session.Database("todo_micro_jois").Collection("courses")

	curr, err := collection.Find(context.TODO(), bson.M{}, findOptions)
	count, err := collection.CountDocuments(context.TODO(), bson.M{}, countOptions)
	if err != nil {
		return &coursepb.ListCourseRes{}, err
	}

	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		err := curr.Decode(data)
		if err != nil {
			log.Fatalln(err)
		}
		courses = append(courses, &coursepb.Course{
			Id:                   data.ID.Hex(),
			Name:                 data.Name,
			Semester:                data.Semester,
		})
	}

	return &coursepb.ListCourseRes{Course: courses, Total: int32(count)}, err
}

func (s *server) CreateCourse(ctx context.Context, req *coursepb.CreateCourseReq) (*coursepb.CreateCourseRes, error) {
	course := &coursepb.Course{
		//Id: 				pmongo.NewObjectId(primitive.NewObjectID()),
		Name:                 req.GetCourse().Name,
		Semester:                req.GetCourse().Semester,
	}
	session, err := db.GetMongoSession()
	if err != nil {
		return &coursepb.CreateCourseRes{}, err
	}
	collection := session.Database("todo_micro_jois").Collection("courses")

	result, err := collection.InsertOne(context.TODO(), course)
	log.Println("result insert", result)

	oid := result.InsertedID.(primitive.ObjectID)
	course.Id = oid.Hex()

	return &coursepb.CreateCourseRes{Course: course}, err
}

func (s *server) ReadCourse(ctx context.Context, req *coursepb.ReadCourseReq) (*coursepb.ReadCourseRes, error) {
	courseID := req.GetId()
	id, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return &coursepb.ReadCourseRes{Course: &coursepb.Course{}}, err
	}

	course  := &Course{}
	session, err := db.GetMongoSession()
	if err != nil {
		return &coursepb.ReadCourseRes{Course: &coursepb.Course{}}, err
	}
	collection := session.Database("todo_micro_jois").Collection("courses")
	query := bson.M{"_id": id}

	err = collection.FindOne(context.TODO(), query).Decode(course)

	return &coursepb.ReadCourseRes{Course: &coursepb.Course{
		Id:                   course.ID.Hex(),
		Name:                 course.Name,
		Semester:                course.Semester,
	}}, err
}

func (s *server) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseReq) (*coursepb.UpdateCourseRes, error) {
	course := req.GetCourse()

	id, err := primitive.ObjectIDFromHex(course.GetId())
	if err != nil {
		return &coursepb.UpdateCourseRes{Course: &coursepb.Course{}}, err
	}

	update := bson.M{
		"name": course.GetName(),
		"semester": course.GetSemester(),
	}

	filter := bson.M{ "_id": id}
	session, err := db.GetMongoSession()
	if err != nil {
		return &coursepb.UpdateCourseRes{Course: &coursepb.Course{}}, err
	}
	collection := session.Database("todo_micro_jois").Collection("courses")

	result := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
	decoded := Course{}
	err = result.Decode(&decoded)
	if err != nil {
		return &coursepb.UpdateCourseRes{Course: &coursepb.Course{}}, err
	}

	return &coursepb.UpdateCourseRes{Course: &coursepb.Course{
		Id: decoded.ID.Hex(),
		Name: decoded.Name,
		Semester: decoded.Semester,
	}}, nil
}

func (s *server) DeleteCourse(ctx context.Context, req *coursepb.DeleteCourseReq) (*coursepb.DeleteCourseRes, error) {
	courseID := req.GetId()
	id, err  := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return &coursepb.DeleteCourseRes{Success:false}, err
	}
	session, err := db.GetMongoSession()
	if err != nil {
		return &coursepb.DeleteCourseRes{Success:false}, err
	}
	collection := session.Database("todo_micro_jois").Collection("courses")

	result, err := collection.DeleteOne(context.TODO(), bson.M{ "_id": id })
	log.Println("result delete", result)

	return &coursepb.DeleteCourseRes{Success: true}, err
}
