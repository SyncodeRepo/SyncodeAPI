package classes

import (
	"database/sql"
	"encoding/json"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/database"
	"github.com/aws/aws-lambda-go/events"
)

type Class struct {
	ID               string
	ClassName        string
	SubjectName      string
	TeacherID        string
	ClassDescription string
}

func HandleGetClass(classID string) (events.APIGatewayProxyResponse, error) {
	return GetClassById(classID)
}

func GetClassById(classId string) (events.APIGatewayProxyResponse, error) {
	var class Class
	err := database.Db.QueryRow("SELECT class_id, class_name, subject_name, teacher_id, class_description FROM Class WHERE class_id = ?", classId).Scan(&class.ID, &class.ClassName, &class.SubjectName, &class.TeacherID, &class.ClassDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			return events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				},
				StatusCode: 404,
				Body:       "Class not found",
			}, nil
		}
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error retrieving class: " + err.Error(),
		}, err
	}
	classJson, err := json.Marshal(class)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error marshalling class data: " + err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 200,
		Body:       string(classJson),
	}, nil
}

func HandleAddClass(requestBody string) (events.APIGatewayProxyResponse, error) {
	var class Class
	err := json.Unmarshal([]byte(requestBody), &class)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 400,
			Body:       "Invalid request body",
		}, err
	}

	stmt, err := database.Db.Prepare("INSERT INTO Class (class_name, subject_name, teacher_id, class_description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error preparing statement",
		}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(class.ClassName, class.SubjectName, class.TeacherID, class.ClassDescription)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error executing statement",
		}, err
	}

	classID, err := result.LastInsertId()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error retrieving last insert ID",
		}, err
	}

	stmt, err = database.Db.Prepare("INSERT INTO teacher_classes (teacher_id, class_id) VALUES (?, ?)")
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error preparing statement for teacher_classes",
		}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(class.TeacherID, classID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Error executing statement for teacher_classes",
		}, err
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 200,
		Body:       "Class added successfully",
	}, nil
}
