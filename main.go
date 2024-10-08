package main

import (
	"database/sql"
	"fmt"
	"html/template"
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

	// TODO: fix templating part
	// Parse and execute the HTML template with the todos data
	tmpl := template.Must(template.New("index").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	 <title>Todo List</title>
	</head>
	<body>
	 <h1>Todo List</h1>
	 <form action="/create" method="POST">
	  <input type="text" name="DESCRIPTION" placeholder="New Todo" required>
	  <input type="checkbox" id="IsCompleted" name="IsCompleted" value="1">
	  <label for="IsCompleted"> Completed</label>
	  <button type="submit">ADD</button>
	 </form>
	 <ul>
	  {{range .}}
	  <li>{{.Description}} <a href="/delete?ID={{.Id}}">Delete</a></li>
	  {{end}}
	 </ul>
	</body>
	</html>
	`))

	tmpl.Execute(w, todos) // Render the template with the list of todos
}

// createHandler handles the creatopn of a new TODO
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// To Extract Description from the form
		DESCRIPTION := r.FormValue("DESCRIPTION")
		//To Determine if Checkbox is Unchecked or Not
		// 1 = Checked
		// 0 = Unchecked
		var IsCompleted int
		if r.FormValue("IsCompleted") == "1" {
			IsCompleted = 1
		} else {
			IsCompleted = 0
		}
		_, err := DB.Exec("INSERT INTO Todos(DESCRIPTION,IsCompleted) VALUES (?,?)", DESCRIPTION, IsCompleted)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

// DeleteHandler handles the Deletion Process
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("ID")
	_, err := DB.Exec("DELETE FROM Todos WHERE ID = ?", ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// TODO: handler for GET Todo

func main() {
	initDB()
	defer DB.Close()

	// Router the handlers for each URL path
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
