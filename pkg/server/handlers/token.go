package handlers

import (
	"encoding/json"
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

	pwd := os.Getenv("TODO_PASSWORD")
	if pwd != password.Password {
		SendError(writer, "неверный пароль", http.StatusUnauthorized)
		return
	}

	iat := time.Now()
	exp := iat.Add(8 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": password.Password,
		"exp":           exp.Unix(),
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
		expectedPassword := os.Getenv("TODO_PASSWORD")
		if expectedPassword == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(expectedPassword), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Invalid or missing token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error": "Invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		passwordHash, ok := claims["password_hash"].(string)
		if !ok || passwordHash != expectedPassword {
			http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
