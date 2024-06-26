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
				StatusCode: 404,
				Body:       "Class not found",
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error retrieving class: " + err.Error(),
		}, err
	}
	classJson, err := json.Marshal(class)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error marshalling class data: " + err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(classJson),
	}, nil
}
