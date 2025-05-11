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
	id INTEGER PRIMARY KEY AUTOINCREMENT,
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

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		panic("Не удалось открыть/создать БД: " + err.Error())
	}

	if install {
		fmt.Println("Базы не было, выполняем скрипт")
		_, err2 := db.Exec(schema)
		if err2 != nil {
			panic("не удалось выполнить скрипт" + err2.Error())
		}
	} else {
		fmt.Println("База уже существует. Продолжаем")
	}

	return nil
}
