package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Db *sql.DB

func Init(databaseName, user, password string) error {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		user, password, databaseName)
	var err error
	Db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %v", err)
	}
	if err := Db.Ping(); err != nil {
		return fmt.Errorf("не удалось пингануть БД: %v", err)
	}

	fmt.Println("Подключение к БД успешно")
	return nil
}
