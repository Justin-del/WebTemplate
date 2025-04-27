package sessions

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTableIfNotExists(){
	database,_ := sql.Open("sqlite3", "file:Database.sqlite");
	database.Exec("create table if not exists sessions(id TEXT PRIMARY KEY)");
}