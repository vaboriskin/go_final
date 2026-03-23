package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустой repeat")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("Неправильный формат")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неверный день")
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
	default:
		return "", errors.New("Неверный формат")
	}
	return date.Format(dateFormat), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowString == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowString)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	fmt.Fprint(w, next)
}
