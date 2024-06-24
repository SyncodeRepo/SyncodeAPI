package users

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/SyncodeRepo/SyncodeAPI.git/endpoints/database"
	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Role      int
	Email     string
}

func HandleGetUsers() events.APIGatewayProxyResponse {
	// Call getUsers to fetch data from the database
	users, err := getUsers()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}
	}

	// Marshal the data into JSON
	jsonData, err := json.Marshal(users)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}
	}

	// Return the JSON data with a 200 OK status
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}
}

func HandleGetUser(ID string) events.APIGatewayProxyResponse {
	user, err := getUserByID(ID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}
	}
	if user == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "User not found",
		}
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}
}

func HandlePostUser(requestBody string) events.APIGatewayProxyResponse {
	var user User

	err := json.Unmarshal(([]byte(requestBody)), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 400,
			Body:       "Error parsing request body",
		}
	}

	stmt, err := database.Db.Prepare("INSERT INTO User (ID, FirstName, LastName, Role, Email) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",                // Adjust this as per your requirements
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS", // Include other methods as needed
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Failed to prepare database statement: " + err.Error(),
		}
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.FirstName, user.LastName, user.Role, user.Email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			},
			StatusCode: 500,
			Body:       "Failed to execute database statement: " + err.Error(),
		}
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		StatusCode: 201,
		Body:       "User successfully created",
	}
}

func getUsers() ([]User, error) {
	rows, err := database.Db.Query("SELECT ID, FirstName, LastName, Role, Email FROM User")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return users, nil
}

func getUserByID(ID string) (*User, error) {
	var user User
	err := database.Db.QueryRow("SELECT ID, FirstName, LastName, Role, Email FROM User WHERE ID = ?", ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err
	}
	return &user, nil
}
