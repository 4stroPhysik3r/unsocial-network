package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func SaveMessageToDB(db *sql.DB, text, sender, receiver string, timeStamp time.Time) error {
	if (text == "") || (sender == "") || (receiver == "") {
		return errors.New("SaveMsgToDB: Missing required fields")
	}
	_, err := db.Exec("INSERT INTO messages (text, timeStamp, sender, receiver) VALUES (?, ?, ?, ?)",
		text, timeStamp, sender, receiver)
	if err != nil {
		log.Println("Error inserting message into the database:", err)
		return err
	}

	return nil
}
