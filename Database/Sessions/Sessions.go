package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTableIfNotExists() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("create table if not exists sessions(id TEXT PRIMARY KEY, absolute_time_out TEXT NOT NULL, idle_time_out TEXT NOT NULL, user_id TEXT NOT NULL, FOREIGN KEY(user_id) REFERENCES users(id))")

}

func CreateASession(user_id string) string {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	id := make([]byte, 32)

	rand.Read(id)

	encoded_id := hex.EncodeToString(id)

	//The idle timeout is common to be 15 minutes for low risk applications and 5 minutes for high risk applications. See https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#session-expiration Also, you might want to lower the absolute timeout for high risk applications.
	database.Exec("insert into sessions (id,absolute_time_out, idle_time_out, user_id) values (?,datetime('now','+24 hour'), datetime('now','+15 minute'),?)", encoded_id, user_id)

	return encoded_id
}

func DoesSessionExistsInDatabase(session_id string) bool {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")

	//Before that, delete all sessions that are expired.
	database.Exec("delete from sessions where datetime('now')>absolute_time_out OR datetime('now')>idle_time_out")

	var result bool
	database.QueryRow("select exists(select 1 FROM sessions where id=?)", session_id).Scan(&result)

	return result
}

func DeleteSession(session_id string) {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")

	database.Exec("delete from sessions where id=?", session_id)
}
