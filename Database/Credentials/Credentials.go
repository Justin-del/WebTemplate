package credentials

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTableIfNotExists() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("Create table if not exists credentials(id BLOB PRIMARY KEY, public_key BLOB NOT NULL, user_id TEXT NOT NULL, signature_counter INTEGER NOT NULL,  FOREIGN KEY(user_id) REFERENCES users(id))")
}

func GetPublicKeyAndSignatureCounter(credential_id []byte, user_id []byte) ([]byte, uint32) {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")

	var public_key []byte
	var signature_counter uint32
	database.QueryRow("select public_key, signature_counter from credentials where user_id=? AND id=?", string(user_id), credential_id).Scan(&public_key, &signature_counter)

	return public_key, signature_counter
}

func UpdateSignatureCounter(credential_id []byte, signature_counter uint32) {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("update credentials set signature_counter=? where id=?", signature_counter, credential_id)
}
