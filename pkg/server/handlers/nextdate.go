package handlers

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
		return "", nil
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
		if len(split) == 1 {
			return "", errors.New("не верный формат repeat")
		}
		days, err := strconv.Atoi(split[1])
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", nil
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", errors.New("не верный формат repeat")
	}

	return date.Format(format), nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDayHandler(response http.ResponseWriter, request *http.Request) {
	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s\nQuery: %s\n",
		request.Method, request.Host, request.URL.Path, request.URL.Query())

	fmt.Println(s)

	query := request.URL.Query()
	dstart := query.Get("date")
	repeat := query.Get("repeat")
	nowFromQuerry := query.Get("now")
	var now time.Time
	var err error
	if nowFromQuerry == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(format, nowFromQuerry)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			now = time.Now()
		}
	}

	date, err := NextDate(now, dstart, repeat)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}

	response.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(response, "%s", date)
}
