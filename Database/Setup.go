package Database

import (
	AuthenticationChallenges "TodoApp/Database/AuthenticationChallenges"
	Credentials "TodoApp/Database/Credentials"
	Sessions "TodoApp/Database/Sessions"
	Todos "TodoApp/Database/Todos"
	Users "TodoApp/Database/Users"
)

func CreateTablesIfNotExists() {
	AuthenticationChallenges.CreateTableIfNotExists()
	Users.CreateTableIfNotExists()
	Credentials.CreateTableIfNotExists()
	Sessions.CreateTableIfNotExists()
	Todos.CreateTableIfNotExists()
}
