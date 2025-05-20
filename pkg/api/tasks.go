package api

import (
	"net/http"

	"github.com/sergeyanosov/go_todo_project/pkg/db"
	"github.com/sergeyanosov/go_todo_project/pkg/server/handlers"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		// ...
		handlers.SendError(w, err, http.StatusInternalServerError)
		return
	}
	handlers.WriteJson(w, TasksResp{
		Tasks: tasks,
	})
}
