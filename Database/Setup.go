package Database

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	Credentials "WebTemplate/Database/Credentials"
	Users "WebTemplate/Database/Users"
)

func CreateTablesIfNotExists() {
	AuthenticationChallenges.CreateTableIfNotExists()
	Users.CreateTableIfNotExists()
	Credentials.CreateTableIfNotExists()
}
