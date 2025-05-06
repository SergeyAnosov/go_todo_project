package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Run() error {
	fmt.Println("Запускаем сервер")
	fmt.Println("Считываем файл .env")
	godotenv.Load()

	http.Handle("/", http.FileServer(http.Dir("/web")))

	addr := os.Getenv("SERVER_ADDRESS")
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s:%d", addr, port)
	return http.ListenAndServe(url, nil)
}
