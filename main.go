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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	webDir := currentDir + "/web"

	err = db.Init(os.Getenv("TODO_DB"))
	if err != nil {
		fmt.Println(err)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		fmt.Println("password: ", pass)
	}

	portEnv := os.Getenv("TODO_PORT")
	if portEnv == "" {
		portEnv = "8080"
	}
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = server.Run(webDir, os.Getenv("TODO_ADDRESS"), port)
	if err != nil {
		fmt.Println(err)
		return
	}
}
