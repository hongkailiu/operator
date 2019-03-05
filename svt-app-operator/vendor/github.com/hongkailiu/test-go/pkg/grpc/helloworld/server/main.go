//https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/hongkailiu/test-go/pkg/probuf/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) GetPerson(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	p := pb.Person{
		Id:    in.Id,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.Person_HOME},
		},
	}

	return &pb.HelloReply{Person: &p}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
