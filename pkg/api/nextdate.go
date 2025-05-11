package api

import (
	"errors"
	"time"
)

var dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat не задан")
	}

	_, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("не удалось распарсить dstart")
	}

	return "", nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
