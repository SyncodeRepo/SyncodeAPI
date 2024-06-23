package students

import (
	"encoding/json"
	"log"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/database"
	"github.com/aws/aws-lambda-go/events"
)

type Class struct {
	ClassID          int    `json:"class_id"`
	ClassName        string `json:"class_name"`
	SubjectName      string `json:"subject_name"`
	TeacherID        int    `json:"teacher_id"`
	ClassDescription string `json:"class_description"`
}

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

func HandleGetStudentClasses(studentID string) events.APIGatewayProxyResponse {
	rows, err := database.Db.Query(`
		SELECT c.class_id, c.class_name, c.subject_name, c.teacher_id, c.class_description
		FROM class c
		JOIN student_classes sc ON c.class_id = sc.class_id
		WHERE sc.student_id = ?`, studentID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to execute query: " + err.Error(),
		}
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ClassID, &class.ClassName, &class.SubjectName, &class.TeacherID, &class.ClassDescription); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error scanning row: " + err.Error(),
			}
		}
		classes = append(classes, class)
	}
	if err := rows.Err(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error iterating over rows: " + err.Error(),
		}
	}

	jsonData, err := json.Marshal(classes)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error marshaling JSON: " + err.Error(),
		}
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 200,
		Body:       string(jsonData),
	}
}
