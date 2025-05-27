package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type PasswordRequest struct {
	Password string `json:"password"`
}

func SigninHandler(writer http.ResponseWriter, request *http.Request) {
	var password PasswordRequest
	body, err := io.ReadAll(request.Body)
	if err != nil {
		SendError(writer, "ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	err = json.Unmarshal(body, &password)
	if err != nil {
		SendError(writer, "ошибка десереализации", http.StatusBadRequest)
		return
	}

	pwd := os.Getenv("PASSWORD")
	if pwd != password.Password {
		SendError(writer, "неверный пароль", http.StatusUnauthorized)
		return
	}

	//todo логику создания токена

	WriteJson(writer, map[string]string{
		"token": "new token",
	})
}
