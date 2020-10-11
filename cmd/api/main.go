package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Todo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var db *sql.DB

func GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Query MySQL to get all the todos from the todos table.
	results, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Println(err)
		w.Write([]byte("We got an error"))
		return
	}
	defer results.Close()

	// Make a variable to hold all of our todos.
	todos := []Todo{}

	// loop through the results aka the row with .Next()
	for results.Next() {
		// crate a todo variable to refernce the Todo struct
		todo := Todo{}
		err := results.Scan(&todo.ID, &todo.Name, &todo.Title, &todo.Description)
		if err != nil {
			log.Println(err)
			w.Write([]byte("We got an error"))
			return
		}
		// Place new todo inside our slice of todos.
		todos = append(todos, todo)
	}

	// Now that we have our slice of todos, let's write
	// back a response of all the todos as JSON.
	jsonBytes, err := json.Marshal(todos)
	if err != nil {
		log.Println(err)
		w.Write([]byte("We got an error"))
		return
	}

	// Set Content-Type header to signify we're sending back
	// JSON, and write out the JSON.
	fmt.Println("json", string(jsonBytes))
	w.Write(jsonBytes)
}

func SingleTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content Type", "application/json; charset=utf-8")
	// fmt.Fprintf(w, "hello, %s!\n", ps.ByName("id"))
	params := ps.ByName("id")

	todo := Todo{}
	selectedTodo := `SELECT id, name, title, description FROM todos WHERE id=?`
	row := db.QueryRow(selectedTodo, params)

	err := row.Scan(&todo.ID, &todo.Name, &todo.Title, &todo.Description)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("we got an error"))
		return
	}
	jsonData, err := json.Marshal(todo)

	if err != nil {
		log.Println(err)
		w.Write([]byte("we got an error"))
		return
	}
	// string converts the bytes to the string representiation aka value
	fmt.Println("json", string(jsonData))
	w.Write(jsonData)
}

func CreateTodo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody Todo
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println(err)
		w.Write([]byte("bad request"))
		return
	}

	// log out the body
	fmt.Printf("body %+v", reqBody)

	insertStatment := `INSERT INTO todos (name, title, description) VALUES (?, ?, ?)`
	_, err := db.Exec(insertStatment, reqBody.Name, reqBody.Title, reqBody.Description)

	jsonData, err := json.Marshal(reqBody)
	fmt.Println("row", jsonData)

	if err != nil {
		log.Println(err)
		w.Write([]byte("We got an error"))
		return
	}
	// Set Content-Type header to signify we're sending back
	// JSON, and write out the JSON.
	fmt.Println("json", string(jsonData))
	// res.json equzlivelnt
	w.Write(jsonData)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	deleteTodo := `DELETE FROM todos WHERE id=?`
	// store a varible with the struct
	// get the param value of the oroute
	params := ps.ByName("id")
	// satrt the qury aosing in the deletetodo and the params
	// be reeutrned in the row
	_, err := db.Exec(deleteTodo, params)
	if err != nil {
		log.Println(err)
		w.Write([]byte("we got an error"))
		return
	}
	w.Write([]byte("deleted"))
}

func UpdateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	updateTodo := `UPDATE todos SET name=?, title=?, description =? WHERE id=?`
	params := ps.ByName("id")

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println(err)
		w.Write([]byte("error"))
		return
	}

	_, err := db.Exec(updateTodo, todo.Name, todo.Title, todo.Description, params)

	if err != nil {
		log.Println(err)
		w.Write([]byte("we have an error"))
	}
	w.Write([]byte("updated"))
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:password@/todo")

	if err != nil {
		panic(err)
	}

	defer db.Close()
	// Ping check were connected
	db.Ping()
	fmt.Println("connected!")

	router := httprouter.New()
	router.GET("/api/todos", GetAll)
	router.GET("/api/todo/:id", SingleTodo)
	router.POST("/api/create/todos", CreateTodo)
	router.PUT("/api/update/:id", UpdateTodo)
	router.DELETE("/api/delete/:id", DeleteTodo)

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
