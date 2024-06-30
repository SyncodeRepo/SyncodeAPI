package students

import (
	"database/sql"
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
		SELECT c.class_id, c.class_name, c.subject_name, c.teacher_id, c.class_description, 
		       a.assignment_id, a.assignment_name, sa.file_url, sa.submission_date, sa.grade
		FROM class c
		LEFT JOIN student_classes sc ON c.class_id = sc.class_id
		LEFT JOIN assignments a ON c.class_id = a.class_id
		LEFT JOIN student_assignments sa ON a.assignment_id = sa.assignment_id AND sa.student_id = ?
		WHERE sc.student_id = ?
		ORDER BY c.class_id, a.assignment_id`, studentID, studentID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to execute query: " + err.Error(),
		}
	}
	defer rows.Close()

	type Assignment struct {
		AssignmentID   int     `json:"assignment_id"`
		AssignmentName string  `json:"assignment_name"`
		FileURL        string  `json:"file_url"`
		SubmissionDate string  `json:"submission_date"`
		Grade          float64 `json:"grade"`
	}

	type ClassWithAssignments struct {
		ClassID          int          `json:"class_id"`
		ClassName        string       `json:"class_name"`
		SubjectName      string       `json:"subject_name"`
		TeacherID        string       `json:"teacher_id"`
		ClassDescription string       `json:"class_description"`
		Assignments      []Assignment `json:"assignments"`
	}

	var classesMap = make(map[int]*ClassWithAssignments)
	for rows.Next() {
		var (
			classID          int
			className        string
			subjectName      string
			teacherID        string
			classDescription string
			assignmentID     sql.NullInt64
			assignmentName   sql.NullString
			fileURL          sql.NullString
			submissionDate   sql.NullString
			grade            sql.NullFloat64
		)
		if err := rows.Scan(&classID, &className, &subjectName, &teacherID, &classDescription, &assignmentID, &assignmentName, &fileURL, &submissionDate, &grade); err != nil {
			log.Printf("Error executing query: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error scanning row: " + err.Error(),
			}
		}

		class, exists := classesMap[classID]
		if !exists {
			class = &ClassWithAssignments{
				ClassID:          classID,
				ClassName:        className,
				SubjectName:      subjectName,
				TeacherID:        teacherID,
				ClassDescription: classDescription,
				Assignments:      []Assignment{},
			}
			classesMap[classID] = class
		}

		if assignmentID.Valid && assignmentName.Valid {
			class.Assignments = append(class.Assignments, Assignment{
				AssignmentID:   int(assignmentID.Int64),
				AssignmentName: assignmentName.String,
				FileURL:        fileURL.String,
				SubmissionDate: submissionDate.String,
				Grade:          grade.Float64,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error iterating over rows: " + err.Error(),
		}
	}

	var classes []ClassWithAssignments
	for _, class := range classesMap {
		classes = append(classes, *class)
	}

	jsonData, err := json.Marshal(classes)
	if err != nil {
		log.Printf("Error executing query: %v", err)
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

func HandleAddStudentToClass(studentID string, classID int) events.APIGatewayProxyResponse {
	_, err := database.Db.Exec(`
		INSERT INTO student_classes (student_id, class_id)
		VALUES (?, ?)`, studentID, classID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Failed to add student to class: " + err.Error(),
		}
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 200,
		Body:       "Student added to class successfully",
	}
}
