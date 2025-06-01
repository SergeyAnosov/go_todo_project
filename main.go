package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sergeyanosov/go_todo_project/pkg/server"

	"github.com/sergeyanosov/go_todo_project/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Старт приложеия")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	webDir := currentDir + "/web"

	err = db.Init(os.Getenv("DB_FILE"))
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}

	err1 := server.Run(webDir, os.Getenv("SERVER_ADDRESS"), port)
	if err1 != nil {
		panic(err1)
	}
}
