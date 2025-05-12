package users

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTableIfNotExists() {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("create table if not exists users(id TEXT PRIMARY KEY, username TEXT UNIQUE NOT NULL)")
}

func DoesUserNameExistsInDatabase(user_name string) bool {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	var username_exists bool

	err := database.QueryRow("select exists(select 1 from users where username=?)", user_name).Scan(&username_exists)

	fmt.Println(err)
	return username_exists
}

func DeleteUser(userId string) {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	database.Exec("pragma foreign_keys=on")
	database.Exec("delete from users where id = ?", userId)
}

/*Returns true if the operation is succesful and false if the operation is not succesful.*/
func AddUserIntoDatabaseWithCredentials(user_id string, user_name string, credential_id []byte, public_key []byte, signature_counter uint32) bool {
	database, _ := sql.Open("sqlite3", "file:Database.sqlite")
	transaction, err := database.Begin()

	if err != nil {
		transaction.Rollback()
		return false
	}

	_, err2 := transaction.Exec("insert into users values(?,?)", user_id, user_name)

	if err2 != nil {
		transaction.Rollback()
		return false
	}

	_, err3 := transaction.Exec("insert into credentials values(?,?,?,?)", credential_id, public_key, user_id, signature_counter)

	if err3 != nil {
		transaction.Rollback()
		return false
	}

	transaction.Commit()
	return true
}
