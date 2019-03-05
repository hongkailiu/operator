//https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/hongkailiu/test-go/pkg/probuf/gen/proto"
	"google.golang.org/grpc"
)

const (
	address   = "localhost:50051"
	defaultID = int32(23)
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)

	// Contact the server and print out its response.
	id := defaultID
	if len(os.Args) > 1 {
		id64, err := strconv.ParseInt(os.Args[1], 10, 32)
		if err != nil {
			log.Fatalf("could not ParseInt: %v", err)
		}
		id = int32(id64)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetPerson(ctx, &pb.HelloRequest{Id: id})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Person: %v", r.Person)
}
