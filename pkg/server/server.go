package server

import (
	"fmt"
	"net/http"
)

func Run(webDir, url string, port int) error {

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	schema := fmt.Sprintf("%s:%d", url, port)
	fmt.Println("Запускаем сервер")
	return http.ListenAndServe(schema, nil)
}
