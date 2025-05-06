package todos

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Id   int
	Todo string
}

func CreateTableIfNotExists() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("create table if not exists todos(id INTEGER PRIMARY KEY, todo TEXT NOT NULL, user_id TEXT NOT NULL, FOREIGN KEY(user_id) references users(id))")
}

/*
Returns the id of the todo that was added.
*/
func AddTodo(todo string, session_id string) int {
	var todoId int
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.QueryRow("insert into todos (todo,user_id) values (?, (select user_id from sessions where sessions.id=?)) RETURNING id", todo, session_id).Scan(&todoId)

	return todoId
}

func GetTodos(session_id string) []Todo {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	rows, _ := database.Query("select id, todo from todos where todos.user_id =  (select user_id from sessions where sessions.id=?)", session_id)
	var todos []Todo

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.Id, &todo.Todo)
		todos = append(todos, todo)
	}

	return todos
}
