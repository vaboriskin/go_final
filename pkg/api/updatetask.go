package api

import (
	"encoding/json"
	"go_final/pkg/db"
	"net/http"
)

func updatetaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан id"}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "пустой title"}, http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{}, http.StatusOK)
}
