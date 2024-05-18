package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var db *sql.DB

func init() {
	var err error
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database)
	db, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Set maximum number of idle connections in the pool.
	db.SetMaxIdleConns(10)

	// Set maximum number of open connections to the database.
	db.SetMaxOpenConns(100)

	// Set the maximum lifetime of a connection.
	db.SetConnMaxLifetime(time.Hour)
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

var (
	server   = "syncode-db.mysql.database.azure.com"
	port     = 3306
	user     = "katamyra"
	password = os.Getenv("DB_PASS")
	database = "syncode"
)

func getUsers() ([]User, error) {
	// Use the global db connection pool
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return users, nil
}
