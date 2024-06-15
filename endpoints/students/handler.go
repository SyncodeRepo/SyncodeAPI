package students

import (
	"encoding/json"
	"log"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/database"
	"github.com/aws/aws-lambda-go/events"
)

type Student struct {
	ID        string `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
}

func HandlePostStudents(requestBody string) events.APIGatewayProxyResponse {
	var student Student
	err := json.Unmarshal([]byte(requestBody), &student)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 400,
			Body:       "Invalid request body",
		}
	}
	stmt, err := database.Db.Prepare("INSERT INTO Student (ID, FirstName, LastName, Email) VALUES (?, ?, ?, ?)")
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error preparing statement",
		}
	}
	defer stmt.Close()
	_, err = stmt.Exec(student.ID, student.FirstName, student.LastName, student.Email)
	if err != nil {
		// Log the values of the variables
		log.Printf("Attempting to execute with ID: %s, FirstName: %s, LastName: %s, Email: %s", student.ID, student.FirstName, student.LastName, student.Email)
		log.Printf("Error executing statement: %v", err) // Log the error
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error executing statement",
		}
	}
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 200,
		Body:       "Student added successfully",
	}
}
