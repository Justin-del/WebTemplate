package authenticationchallenge

import (
	webauthn "WebTemplate/Utils/WebAuthn"
	"crypto/rand"
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTableIfNotExists() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("Create table if not exists authentication_challenges(id INTEGER PRIMARY KEY, challenge BLOB NOT NULL, expires_at TEXT NOT NULL)")
}

func CreateNewChallenge() webauthn.Challenge {
	challenge := make([]byte, 32) // 256 bits = 32 bytes
	rand.Read(challenge)

	timeout := "+" + strconv.Itoa(webauthn.TimeoutInMinutes) + " minute"

	var id int

	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.QueryRow("insert into authentication_challenges(challenge, expires_at) values (?,datetime('now',?)) RETURNING id", challenge, timeout).Scan(&id)

	return webauthn.Challenge{
		Id:        id,
		Challenge: challenge,
	}
}

/*
Returns the challenge that was deleted
*/
func DeleteChallengeByID(id any) []byte {
	var challenge []byte
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.QueryRow("delete from authentication_challenges where id=? returning challenge", id).Scan(&challenge)
	return challenge
}

func DeleteAnyExpiredChallenges() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("delete from authentication_challenges where datetime('now')>expires_at")
}
