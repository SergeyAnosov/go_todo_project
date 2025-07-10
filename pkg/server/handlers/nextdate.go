package handlers

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var format = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (pgtype.Date, error) {
	if len(repeat) == 0 {
		return pgtype.Date{}, nil
	}
	if strings.Contains(repeat, "w ") || strings.Contains(repeat, "m ") {
		return pgtype.Date{}, errors.New("не поддерживаемый формат")
	}
	date, err := time.Parse(format, dstart)
	if err != nil {
		return pgtype.Date{}, errors.New("не удалось распарсить dstart")
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
			return pgtype.Date{}, errors.New("не верный формат repeat")
		}
		days, err := strconv.Atoi(split[1])
		if err != nil {
			return pgtype.Date{}, err
		}
		if days > 400 {
			return pgtype.Date{}, nil
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return pgtype.Date{}, errors.New("не верный формат repeat")
	}

	pgDate := pgtype.Date{}
	err = pgDate.Scan(date)
	if err != nil {
		return pgtype.Date{}, errors.New("Что то пошло не так")
	}

	fmt.Printf("time.Time: %v\n", now.Format("2006-01-02"))
	fmt.Printf("pgtype.Date: %v\n", pgDate)
	return pgDate, nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDayHandler(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	dstart := query.Get("date")
	repeat := query.Get("repeat")
	nowFromQuery := query.Get("now")
	var now time.Time
	var err error
	if len(nowFromQuery) == 0 {
		now = time.Now()
	} else {
		now, err = time.Parse(format, nowFromQuery)
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
	response.Write([]byte(date))
}
