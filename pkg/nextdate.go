package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var format = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat не задан")
	}
	if strings.Contains(repeat, "w ") || strings.Contains(repeat, "m ") {
		return "", errors.New("не поддерживаемый формат")
	}
	date, err := time.Parse(format, dstart)
	if err != nil {
		return "", errors.New("не удалось распарсить dstart")
	}

	split := strings.Split(repeat, " ")
	switch split[0] {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	case "d":
		days, err := strconv.Atoi(split[1])
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
	}

	return date.Format(format), nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDayHandler(response http.ResponseWriter, request *http.Request) {
	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		request.Method, request.Host, request.URL.Path)
	response.Write([]byte(s))
}
