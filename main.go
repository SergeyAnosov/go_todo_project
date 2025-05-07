package main

import (
	"fmt"
	"os"
	"strconv"

	"go1f/pkg/server"

	"go1f/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := db.Init(os.Getenv("DB_PATH"))
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}

	err1 := server.Run("/web", os.Getenv("SERVER_ADDRESS"), port)
	if err1 != nil {
		fmt.Println("Завершаем работу")
		fmt.Println(err1)
		panic(err1)
	}
}
