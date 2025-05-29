package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"os"
	"time"
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

	iat := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": password.Password,
		"exp":           iat,
	})

	signedToken, err := token.SignedString([]byte(pwd))
	if err != nil {
		SendError(writer, "не удалось получить jwt token", http.StatusInternalServerError)
	}

	WriteJson(writer, map[string]interface{}{
		"token": signedToken,
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, установлен ли пароль
		pass := os.Getenv("PASSWORD")
		if len(pass) > 0 {
			var jwtToken string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			var valid bool

			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return pass, nil
			})

			if err == nil && token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok {
					hash, ok := claims["password_hash"].(string)
					if ok && hash == pass {
						valid = true
					}
				}
			}

			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
