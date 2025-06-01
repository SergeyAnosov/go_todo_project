package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	err = json.Unmarshal(body, &task)
	if err != nil {
		SendError(writer, "ошибка десереализации", http.StatusBadRequest)
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
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if afterNow(now.Truncate(24*time.Hour), t) {
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

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		SendError(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	WriteJson(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, "ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &task)
	if err != nil {
		SendError(w, "ошибка десереализации: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что заголовок не пустой
	if checkTitleIsEmpty(&task) {
		SendError(w, "заголовок не может быть пустым", http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := task.ID
	_, err = db.GetTask(id)
	if err != nil {
		SendError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		SendError(w, "не удалось обновить задачу", http.StatusInternalServerError)
		return
	}
	WriteJson(w, map[string]string{})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "не указан индентификатор", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			SendError(w, "задача не найдена", http.StatusNotFound)
			return
		}
		SendError(w, "ошибка удаления задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJson(w, map[string]string{})
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		SendError(w, "не указан индентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		SendError(w, "задача не найдена: ", http.StatusNotFound)
		return
	}

	if len(task.Repeat) == 0 {
		err := db.DeleteTask(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				SendError(w, "Задача не найдена", http.StatusNotFound)
				return
			}
			SendError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}
		WriteJson(w, map[string]string{})
		return
	}

	date, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateDate(date, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(w, "задача не найдена", http.StatusNotFound)
			return
		}
		SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, map[string]string{})
}
