package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// TODO: struct for Todo
type Todo struct {
	Id          int64
	Description string
	IsCompleted bool
}

// TODO: initialize DB
var DB *sql.DB

func initDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./Todos.db") // Open a connection to the SQlite database file named Todos.db
	if err != nil {
		log.Fatal(err)
	}

	// note: Table already created
	// // SQL statement to create the todos table if it doesn't exist
	// sqlStmt := `
	// CREATE TABLE IF NOT EXISTS todos (
	//  ID INTEGER NOT NULL PRIMARY KEY,
	//  Description TEXT,
	//  IsCompleted BOOLEAN
	// );`

	// _, err = DB.Exec(sqlStmt)
	// if err != nil {
	//  log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt) // Log an error if table creation fails
	// }
}

// TODO: Index Handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Query the database to get all todos
	rows, err := DB.Query("SELECT ID, Description, IsCompleted FROM Todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Description, &todo.IsCompleted); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	return
}

// TODO: handler for GET Todo

func main() {
	initDB()
	defer DB.Close()

	// Router the handlers for each URL path
	http.HandleFunc("/", indexHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
