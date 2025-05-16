package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sergeyanosov/go_todo_project/pkg/api"
	"net/http"
)

func Run(webDir, url string, port int) error {
	r := chi.NewRouter()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(r)
	schema := fmt.Sprintf("%s:%d", url, port)
	fmt.Println("Запускаем сервер")
	return http.ListenAndServe(schema, r)
}
