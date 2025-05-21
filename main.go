package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sergeyanosov/go_todo_project/pkg/server"

	"github.com/sergeyanosov/go_todo_project/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	webDir := currentDir + "/web"

	err = godotenv.Load()
	if err != nil {
		return
	}
	fmt.Printf("Файл .env успешно считан\n")

	err = db.Init(os.Getenv("DB_PATH"))
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}

	err1 := server.Run(webDir, os.Getenv("SERVER_ADDRESS"), port)
	if err1 != nil {
		fmt.Println("Завершаем работу")
		fmt.Println(err1)
		panic(err1)
	}
}
