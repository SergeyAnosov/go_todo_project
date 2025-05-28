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

	pwd := os.Getenv("PASSWORD")
	if pwd != password.Password {
		SendError(writer, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": password.Password, // можно хранить хэш, но в учебном примере — сам пароль
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString(([]byte(pwd)))
	if err != nil {
		SendError(writer, "не удалось получить jwt token", http.StatusInternalServerError)
	}

	WriteJson(writer, map[string]interface{}{
		"token": signedToken,
	})
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		expectedPassword := os.Getenv("PASSWORD")
		if expectedPassword == "" {
			next(writer, request)
			return
		}

		cookie, err := request.Cookie("token")
		if err != nil {
			SendError(writer, "требуется авторизация", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ([]byte)(expectedPassword), nil
		})

		if err != nil || !token.Valid {
			SendError(writer, "invalid or missing token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			SendError(writer, "invalid token claims", http.StatusUnauthorized)
			return
		}

		passwordHash, ok := claims["password_hash"].(string)
		if !ok || passwordHash != expectedPassword {
			SendError(writer, "authentification required", http.StatusUnauthorized)
			return
		}
		next(writer, request)
	}
}

func AsChiMiddleware(fn func(http.HandlerFunc) http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(fn(next.ServeHTTP))
	}
}
