package main

import (
	"strconv"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		switch request.Resource {
		case "/users":
			// Handle GET /users
			return users.HandleGetUsers(), nil
		case "/classes":
			// Handle GET /classes
			return handleGetClasses(), nil
		case "/users/{id}":
			id, ok := request.PathParameters["id"]
			if !ok {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "ID parameter is missing",
				}, nil
			}
			_, err := strconv.Atoi(id)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Invalid ID format",
				}, nil
			}
			userResponse := users.HandleGetUser(id)
			return events.APIGatewayProxyResponse{
				StatusCode: userResponse.StatusCode,
				Body:       userResponse.Body,
			}, nil
		default:
			// Handle unknown resource
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Resource Not Found",
			}, nil
		}
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid Request Method: " + request.HTTPMethod,
		}, nil
	}
}

func handleGetClasses() events.APIGatewayProxyResponse {
	// Your logic to handle GET /classes
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Handle GET /classes",
	}
}

func main() {
	lambda.Start(handler)
}
