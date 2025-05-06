package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

var schema string = `
CREATE TABLE scheduler (
	id INEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX date_index ON scheduler (date);`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	fmt.Println(install)

	if install {
		_, err := db.Exec(schema)
		return err
	}

	db, err = sql.Open("sqlite", dbFile)
	return err
}
