package api

import (
	"net/http"

	"github.com/sergeyanosov/go_todo_project/pkg"
)

func Init() {
	http.HandleFunc("/api/nextdate", pkg.NextDayHandler)
}
