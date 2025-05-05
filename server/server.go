package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handle(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("Получен запрос: %s\n", req.URL.Path)
}

func Start() {
	fmt.Println("Запускаем сервер")
	fmt.Println("Считываем файл .env")
	godotenv.Load()

	http.Handle("/", http.FileServer(http.Dir("web")))
	//http.HandleFunc(`/`, handle)

	addr := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("TODO_PORT")

	fmt.Printf("port: %s\n", port)

	host_and_Port := fmt.Sprintf("%s:%s", addr, port)
	fmt.Println("host_and_Port:", host_and_Port)
	err := http.ListenAndServe(host_and_Port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}
