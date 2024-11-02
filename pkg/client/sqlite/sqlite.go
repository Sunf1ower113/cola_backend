package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func NewDB(driver, name string) (db *sql.DB, err error) {
	db, err = sql.Open(driver, name)
	if err != nil {
		log.Println("error create db")
		return
	}
	log.Println("create db: OK")
	err = createTable(db)
	if err != nil {
		log.Println("error create tables")
		return
	}
	return
}
func createTable(db *sql.DB) error {
	var query []string
	users := `
CREATE TABLE IF NOT EXISTS users(
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT DEFAULT '',
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	phone_number TEXT DEFAULT '',
	birth_date DATE DEFAULT '',
	points INTEGER DEFAULT 0,
	role TEXT CHECK (role IN ('admin', 'user')) DEFAULT 'user'
)
`
	recycle_boxes := `
CREATE TABLE IF NOT EXISTS recycle_boxes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
    address TEXT NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 10,
    count INTEGER NOT NULL DEFAULT 0

)
`
	query = append(query, users, recycle_boxes)
	for _, v := range query {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func migrate(db *sql.DB) error {
	var query []string
	users := ``
	query = append(query, users)
	for _, v := range query {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
