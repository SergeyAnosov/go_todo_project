package api

import (
	"github.com/sergeyanosov/go_todo_project/pkg/server/handlers"
)

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", handlers.NextDayHandler)
	r.Post("/api/task", handlers.AddTaskHandle)
	r.Get("/api/tasks")

}
