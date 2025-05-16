package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sergeyanosov/go_todo_project/pkg/db"
	"io"
	"net/http"
	"time"
)

func AddTaskHandle(writer http.ResponseWriter, request *http.Request) {
	var task db.Task

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	fmt.Printf("body: %s\n", string(body))

	err = json.Unmarshal(body, &task)
	if err != nil {
		fmt.Printf("Ошибка десереализации:%s\n", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if checkTitleIsEmpty(&task) {
		err = fmt.Errorf("поле Title не задано")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	lastTaskId, err := db.AddTask(&task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	writeJson(writer, lastTaskId)
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
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		fmt.Println(err)
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if now.After(t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}
	return nil
}
func writeJson(w http.ResponseWriter, data any) {
	out, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
