package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//The idle timeout is common to be 15 to 30 minutes for low risk applications and 2-5 minutes for high risk applications. See https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#session-expiration 
var idleTimeOutInMinutes int = 15

func CreateTableIfNotExists() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("create table if not exists sessions(id TEXT PRIMARY KEY, absolute_time_out TEXT NOT NULL, idle_time_out TEXT NOT NULL, user_id TEXT NOT NULL, FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE)")
}

func CreateASession(userId string) string {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	id := make([]byte, 32)

	rand.Read(id)

	encoded_id := hex.EncodeToString(id)

	//You might want to lower the absolute timeout for high risk applications.
	database.Exec("insert into sessions (id,absolute_time_out, idle_time_out, user_id) values (?,datetime('now','+24 hour'), datetime('now',?),?)", encoded_id, "+"+strconv.Itoa(idleTimeOutInMinutes)+" minutes", userId)

	return encoded_id
}

func DoesSessionExistsInDatabase(sessionId string) bool {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")

	//Before that, delete all sessions that are expired.
	database.Exec("delete from sessions where datetime('now')>absolute_time_out OR datetime('now')>idle_time_out")

	var returnedSessionId string

	//Extend the session by the defined idle timeout in minutes.
	database.QueryRow("update sessions set idle_time_out=datetime('now',?) where id=? RETURNING id", "+"+strconv.Itoa(idleTimeOutInMinutes)+" minutes", sessionId).Scan(&returnedSessionId)

	return returnedSessionId != ""
}

func DeleteSession(sessionId string) {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")

	database.Exec("delete from sessions where id=?", sessionId)
}
