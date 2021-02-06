package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

//Establish connection to the database and create required table if iit does not exist
func Connect(path string) {
	if db != nil {
		log.Fatal("db pointer is not nill")
	}
	var err error
	db, err = sql.Open("sqlite3", path+"TradeX.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;

		CREATE TABLE IF NOT EXISTS Events(
		date TEXT PRIMARY KEY NOT NULL,
		views INT NOT NULL,
		clicks INT NOT NULL,
		cost REAL NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

//Close connection to the database
func Close() {
	db.Close()
	db = nil
}
