package main

import (
	"fmt"
	"os"

	"go1f/pkg/server"

	"go1f/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbFile := os.Getenv("DB_PATH")
	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}

	err1 := server.Run()
	if err1 != nil {
		fmt.Println("Завершаем работу")
		fmt.Println(err1)
		panic(err1)
	}
}
