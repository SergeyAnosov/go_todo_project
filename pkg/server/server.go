package server

import (
	"fmt"
	"github.com/sergeyanosov/go_todo_project/pkg/api"
	"net/http"
)

func Run(webDir, url string, port int) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()
	schema := fmt.Sprintf("%s:%d", url, port)
	fmt.Println("Запускаем сервер")
	return http.ListenAndServe(schema, nil)
}
