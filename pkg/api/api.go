package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/sergeyanosov/go_todo_project/pkg/server/handlers"
)

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", handlers.NextDayHandler)

	// Защищённые эндпоинты
	r.With(handlers.Auth).Post("/api/task", handlers.AddTaskHandle)
	r.With(handlers.Auth).Get("/api/tasks", handlers.TasksHandler)
	r.With(handlers.Auth).Get("/api/task", handlers.GetTaskHandler)
	r.With(handlers.Auth).Put("/api/task", handlers.UpdateTaskHandler)
	r.With(handlers.Auth).Delete("/api/task", handlers.DeleteTaskHandler)
	r.With(handlers.Auth).Post("/api/task/done", handlers.TaskDoneHandler)

	r.Post("/api/signin", handlers.SigninHandler)
}
