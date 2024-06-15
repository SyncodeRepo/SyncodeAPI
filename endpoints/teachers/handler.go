package teachers

import (
	"encoding/json"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/database"
	"github.com/aws/aws-lambda-go/events"
)

type Teacher struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	SchoolName string `json:"schoolName"`
}

func HandlePostTeachers(requestBody string) events.APIGatewayProxyResponse {
	var teacher Teacher
	err := json.Unmarshal([]byte(requestBody), &teacher)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 400,
			Body:       "Invalid request body",
		}
	}
	stmt, err := database.Db.Prepare("INSERT INTO Teacher (ID, FirstName, LastName, Email, SchoolName) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error preparing statement",
		}
	}
	defer stmt.Close()
	_, err = stmt.Exec(teacher.ID, teacher.FirstName, teacher.LastName, teacher.Email, teacher.SchoolName)
	if err != nil {
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
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 201,
		Body:       "Teacher added successfully",
	}
}
