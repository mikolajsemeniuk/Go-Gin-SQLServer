package database

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"

	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var server = "localhost"
var port = 1433
var user = "sa"
var password = "P@ssw0rd"
var database = "db"

var Client *sql.DB

func init() {
	var err error

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	fmt.Println(server, user, password, database)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	fmt.Printf("Connected!\n")
	defer conn.Close()

}
