// https://ops.tips/blog/a-swagger-golang-hello-world/
// main declares the CLI that spins up the server of
// our API.
// It takes some arguments, validates if they're valid
// and match the expected type and then intiialize the
// server.
package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/hongkailiu/test-go/pkg/swagger/swagger/models"
	"github.com/hongkailiu/test-go/pkg/swagger/swagger/restapi"
	"github.com/hongkailiu/test-go/pkg/swagger/swagger/restapi/operations"
)

type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	args = &cliArgs{
		Port: 8080,
	}
)

// getHostnameHandler implements the handler that
// takes a set of parameters as described in swagger.yml
// and then produces a response.
// This response might be an error of a succesfull response.
//	-	In case of failure we create the payload
//		that would indicate the failure.
//	-	In case of success, the payload the indicates
//		the success with the hostname.
func getHostnameHandler(params operations.GetHostnameParams) middleware.Responder {
	payload, err := os.Hostname()

	if err != nil {
		errPayload := &models.Error{
			Code:    500,
			Message: swag.String("failed to retrieve hostname"),
		}

		return operations.
			NewGetHostnameDefault(500).
			WithPayload(errPayload)
	}

	return operations.NewGetHostnameOK().WithPayload(payload)
}

func getUsersHandler(params operations.GetUsersParams) middleware.Responder {
	id := int64(3)
	users := []*models.User{{&id, "hongkliu"}}
	return operations.NewGetUsersOK().WithPayload(users)
}

func getUserUserIDHandler(params operations.GetUserUserIDParams) middleware.Responder {

	if params.UserID == 1 {
		id := int64(1)
		user := models.User{&id, "mike"}
		return operations.NewGetUserUserIDOK().WithPayload(&user)
	}
	return operations.NewGetUserUserIDNotFound()

}

// main performs the main routine of the application:
//	1.	parses the args;
//	2.	analyzes the declaration of the API
//	3.	sets the implementation of the handlers
//	4.	listens on the port we want
func main() {
	arg.MustParse(args)

	// Load the JSON that corresponds to our swagger.yml
	// api definition.
	// This JSON is hardcoded as part of the generated code
	// that go-swagger creates.
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// Load a dummy object that servers as an interface
	// that allows us to implement the API specification.
	api := operations.NewHelloAPI(swaggerSpec)

	// Create the REST api server that will make use of
	// the object that will container our handler implementations.
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// Configure the server port
	server.Port = args.Port

	// Add our handler implementation
	api.GetHostnameHandler = operations.GetHostnameHandlerFunc(
		getHostnameHandler)

	api.GetUsersHandler = operations.GetUsersHandlerFunc(getUsersHandler)
	api.GetUserUserIDHandler = operations.GetUserUserIDHandlerFunc(getUserUserIDHandler)

	// Let it run
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
