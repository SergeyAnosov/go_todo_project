package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sergeyanosov/go_todo_project/pkg/db"
	"io"
	"net/http"
	"strconv"
	"time"
)

func AddTaskHandle(writer http.ResponseWriter, request *http.Request) {
	var task db.Task

	body, err := io.ReadAll(request.Body)
	if err != nil {
		SendError(writer, "ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	fmt.Printf("body: %s\n", string(body))

	err = json.Unmarshal(body, &task)
	if err != nil {
		SendError(writer, "ошибка десереализации: "+err.Error(), http.StatusBadRequest)
		return
	}

	if checkTitleIsEmpty(&task) {
		err = fmt.Errorf("не указан заголовок задачи")
		SendError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		SendError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	lastTaskId, err := db.AddTask(&task)
	if err != nil {
		SendError(writer, err.Error(), http.StatusBadRequest)
	}
	WriteJson(writer, struct {
		ID string `json:"id"`
	}{
		ID: strconv.FormatInt(lastTaskId, 10),
	})
}

func checkTitleIsEmpty(task *db.Task) bool {
	if task.Title == "" {
		return true
	} else {
		return false
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format("20060102")
	}
	t, err := time.Parse("20060102", task.Date)
	fmt.Printf("t: %s\n", t)
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	fmt.Printf("next: %s\n", next)
	fmt.Printf("now: %s\n", now)
	if err != nil {
		return err
	}

	if afterNow(now.Truncate(24*time.Hour), t) {
		fmt.Printf("afternow is %v\n", afterNow(now, t))
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}
	return nil
}
func WriteJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Если произошла ошибка при кодировании JSON — отправляем Internal Server Error
		http.Error(w, `{"error":"Ошибка при сериализации данных"}`, http.StatusInternalServerError)
	}
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		// ...
		SendError(w, "ошибка получения всех задач", http.StatusInternalServerError)
		return
	}
	WriteJson(w, TasksResp{
		Tasks: tasks,
	})
}
