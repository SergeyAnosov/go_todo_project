package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sergeyanosov/go_todo_project/pkg/db"
	"github.com/sergeyanosov/go_todo_project/pkg/server"
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

	err = db.Init(os.Getenv("PG_DB"), os.Getenv("PG_USR"), os.Getenv("PG_PWD"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Db.Close()

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
