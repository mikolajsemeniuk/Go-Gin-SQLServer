package database

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

// Replace with your own connection parameters
var server = "localhost"
var port = 1433
var user = "sa"
var password = "P@ssw0rd"

var Client *sql.DB

func init() {
	var error error

	// Create connection string
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
		server, user, password, port)

	// Create connection pool
	Client, error = sql.Open("sqlserver", connectionString)
	if error != nil {
		log.Fatal("Error creating connection pool: " + error.Error())
	}
	log.Printf("Connected!\n")

	// Close the database connection pool after program executes
	defer Client.Close()

	// SelectVersion()
}

// // Gets and prints SQL Server version
// func SelectVersion() {
// 	// Use background context
// 	ctx := context.Background()

// 	// Ping database to see if it's still alive.
// 	// Important for handling network issues and long queries.
// 	err := Client.PingContext(ctx)
// 	if err != nil {
// 		log.Fatal("Error pinging database: " + err.Error())
// 	}

// 	var result string

// 	// Run query and scan for result
// 	err = Client.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
// 	if err != nil {
// 		log.Fatal("Scan failed:", err.Error())
// 	}
// 	fmt.Printf("%s\n", result)
// }
