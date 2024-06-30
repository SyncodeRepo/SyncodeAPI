package main

import (
	"strconv"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/classes"
	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/students"
	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/teachers"
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
		case "/classes/{id}":
			id, ok := request.PathParameters["id"]
			if !ok {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "ID parameter is missing",
				}, nil
			}
			classResponse, err := classes.HandleGetClass(id)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       "Internal Server Error",
				}, err
			}
			return classResponse, nil
		case "/users/{id}":
			id, ok := request.PathParameters["id"]
			if !ok {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "ID parameter is missing",
				}, nil
			}
			userResponse := users.HandleGetUser(id)
			return events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
					"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				},
				StatusCode: userResponse.StatusCode,
				Body:       userResponse.Body,
			}, nil
		case "/students/{id}/classes":
			id, ok := request.PathParameters["id"]
			if !ok {
				return events.APIGatewayProxyResponse {
					StatusCode: 400,
					Body:       "ID parameter is missing",
				}, nil
			}
			return students.HandleGetStudentClasses(id), nil
		case "/teachers/{id}/classes":
			id, ok := request.PathParameters["id"]
			if !ok {
				return events.APIGatewayProxyResponse {
					StatusCode: 400,
					Body:       "ID parameter is missing",
				}, nil
			}
			return teachers.HandleGetTeacherClasses(id), nil
		default:
			// Handle unknown resource
			return events.APIGatewayProxyResponse {
				StatusCode: 404,
				Body:       "Resource Not Found",
			}, nil
		}
	case "POST":
		switch request.Resource {
		case "/users":
			requestBody := request.Body
			return users.HandlePostUser(requestBody), nil
		case "/students/{id}/classes":
			id, ok := request.PathParameters["id"]
			if !ok {
				return events.APIGatewayProxyResponse {
					StatusCode: 400,
					Body:       "ID parameter is missing",
				}, nil
			}
			classIDStr, ok := request.QueryStringParameters["class_id"]
			if !ok {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "class_id query parameter is missing",
				}, nil
			}
			classID, err := strconv.Atoi(classIDStr)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Invalid class_id format",
				}, nil
			}
			return students.HandleAddStudentToClass(id, classID), nil
		case "/students":
			requestBody := request.Body
			return students.HandlePostStudents(requestBody), nil
		case "/teachers":
			requestBody := request.Body
			return teachers.HandlePostTeachers(requestBody), nil
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Resource not found",
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
