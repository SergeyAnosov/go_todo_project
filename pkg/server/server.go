package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sergeyanosov/go_todo_project/pkg/api"
	"net/http"
)

func Run(webDir, url string, port int) error {
	r := chi.NewRouter()
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	api.Init(r)
	schema := fmt.Sprintf("%s:%d", url, port)
	fmt.Printf("Cервер запущен по адесу %s\n", schema)
	return http.ListenAndServe(schema, r)
}
