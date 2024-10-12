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
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Todo List</title>
    <style>
        /* General Styling */
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f4f4f9;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }

        /* Container */
        .container {
            background-color: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            max-width: 500px;
            width: 100%;
        }

        /* Header */
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 20px;
        }

        /* Form styling */
        form {
            display: flex;
            justify-content: space-between;
            margin-bottom: 20px;
        }

        input[type="text"] {
            width: 70%;
            padding: 10px;
            font-size: 1rem;
            border: 1px solid #ddd;
            border-radius: 5px;
            outline: none;
        }

        button {
            padding: 10px 20px;
            background-color: #28a745;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s;
        }

        button:hover {
            background-color: #218838;
        }

        /* Todo List Styling */
        ul {
            list-style-type: none;
            padding: 0;
        }

        li {
            display: flex;
            justify-content: space-between;
            align-items: center;
            background-color: #f9f9f9;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-bottom: 10px;
        }

        li label {
            flex-grow: 1;
            margin-left: 10px;
        }

        /* Completed task style */
        .completed {
            
            color: #888;
        }

        /* Delete link */
        a {
            color: #dc3545;
            text-decoration: none;
            font-weight: bold;
        }

        a:hover {
            text-decoration: underline;
        }

        /* Checkbox Styling */
        input[type="checkbox"] {
            transform: scale(1.5);
            cursor: pointer;
        }
    </style>
</head>
<body>

    <div class="container">
        <h1>Todo List</h1>

        <form action="/create" method="POST">
            <input type="text" name="DESCRIPTION" placeholder="New Todo" required>
            <button type="submit">Add</button>
        </form>

        <ul>
            {{range .}}
            <li>
                <input type="checkbox" id="isCompleted_{{.Id}}" name="isCompleted" {{if .IsCompleted}}checked{{end}}>
                <label for="isCompleted_{{.Id}}" class="{{if .IsCompleted}}completed{{end}}">
                    {{.Description}}
                </label>
                <a href="/delete?ID={{.Id}}" onclick="return confirm('Are you sure you want to delete this item?')">Delete</a>
            </li>
            {{end}}
        </ul>
    </div>

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

func updateHandler(w http.ResponseWriter, r *http.Request) {

}

// TODO: handler for GET Todo

func main() {
	initDB()
	defer DB.Close()

	// Router the handlers for each URL path
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/update", updateHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
