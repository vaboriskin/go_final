package api

import (
	"go_final/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	search := r.URL.Query().Get("search")
	tasks, err := db.SearchTasks(limit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)

}
