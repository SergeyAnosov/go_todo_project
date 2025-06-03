package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

var schema string = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX date_index ON scheduler (date);`

func Init(dbFile string) error {
	dir := filepath.Dir(dbFile)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку %s: %v", dir, err)
	}

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	var errOpen error
	db, errOpen = sql.Open("sqlite", dbFile)
	if errOpen != nil {
		return fmt.Errorf("не удалось открыть/создать БД: %v", errOpen)
	}
	defer db.Close()

	if install {
		_, errSchema := db.Exec(schema)
		if errSchema != nil {
			return fmt.Errorf("не удалось выполнить скрипт создания схемы: %v", errSchema)
		}
	}
	return nil
}
