package db

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const dateFormat = "20060102"

func AddTask(task *Task) (int64, error) {
	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?,?,?,?)
	`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	query := `
	SELECT id, date, title, comment, repeat
	FROM scheduler
	ORDER BY date
	LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `
	SELECT id,date,title, comment, repeat
	FROM scheduler
	WHERE id = ?
	`

	var t Task

	err := DB.QueryRow(query, id).Scan(
		&t.ID,
		&t.Date,
		&t.Title,
		&t.Comment,
		&t.Repeat,
	)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}
	return &t, nil
}

func UpdateTask(task *Task) error {
	query := `
	UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?
	`
	res, err := DB.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `
	DELETE FROM scheduler WHERE id = ?
	`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(date string, id string) error {
	query := `
	UPDATE scheduler SET date = ? WHERE id = ?
	`
	res, err := DB.Exec(query, date, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func SearchTasks(limit int, search string) ([]*Task, error) {
	var query string
	var args []interface{}
	if t, err := time.Parse("02.01.2006", search); err == nil {
		dateString := t.Format(dateFormat)
		query = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE date = ?
		ORDER BY date
		LIMIT ?
		`
		args = []interface{}{dateString, limit}
	} else if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		query = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date
		LIMIT ?
		`
		args = []interface{}{like, like, limit}
	} else {
		query = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
		`
		args = []interface{}{limit}
	}
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}
