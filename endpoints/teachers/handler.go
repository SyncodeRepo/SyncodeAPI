package teachers

import (
	"encoding/json"
	"log" // Add log package

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/database"
	"github.com/aws/aws-lambda-go/events"
)

type Teacher struct {
	ID         string `json:"ID"`
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Email      string `json:"Email"`
	SchoolName string `json:"SchoolName"`
}

type Class struct {
	ID               string `json:"class_id"`
	ClassName        string `json:"class_name"`
	SubjectName      string `json:"subject_name"`
	TeacherID        string `json:"teacher_id"`
	ClassDescription string `json:"class_description"`
}

func HandlePostTeachers(requestBody string) events.APIGatewayProxyResponse {
	log.Println("Received request body:", requestBody) // Log the request body

	var teacher Teacher
	err := json.Unmarshal([]byte(requestBody), &teacher)
	if err != nil {
		log.Println("Error unmarshalling request body:", err) // Log the error
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
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 201,
		Body:       "Teacher added successfully",
	}
}

func HandleGetTeacherClasses(teacherID string) events.APIGatewayProxyResponse {
	rows, err := database.Db.Query(`
		SELECT c.class_id, c.class_name, c.subject_name, c.teacher_id, c.class_description
		FROM class c
		JOIN teacher_classes tc ON c.class_id = tc.class_id
		WHERE tc.teacher_id = ?`, teacherID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Failed to execute query: " + err.Error(),
		}
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.ClassName, &class.SubjectName, &class.TeacherID, &class.ClassDescription); err != nil {
			return events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				},
				StatusCode: 500,
				Body:       "Error scanning row: " + err.Error(),
			}
		}
		classes = append(classes, class)
	}
	if err := rows.Err(); err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error iterating over rows: " + err.Error(),
		}
	}

	jsonData, err := json.Marshal(classes)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error marshaling JSON: " + err.Error(),
		}
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 200,
		Body:       string(jsonData),
	}
}
