package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	server   = "syncode-db.mysql.database.azure.com"
	port     = 3306
	user     = "katamyra"
	password = os.Getenv("DB_PASS")
	database = "syncode"
)

var Db *sql.DB

func init() {
	var err error
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database)
	Db, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Set maximum number of idle connections in the pool.
	Db.SetMaxIdleConns(10)

	// Set maximum number of open connections to the database.
	Db.SetMaxOpenConns(100)

	// Set the maximum lifetime of a connection.
	Db.SetConnMaxLifetime(time.Hour)
}
