package db

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type Task struct {
	ID      int         `json:"id"`
	Date    pgtype.Date `json:"date"`
	Title   string      `json:"title"`
	Comment string      `json:"comment"`
	Repeat  string      `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := "INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task = make([]*Task, 0)
	rows, err := Db.Query("SELECT id, date, title, comment, repeat FROM scheduler LIMIT ?", limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return []*Task{}, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	err := Db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id= :id", sql.Named("id", id)).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`

	res, err := Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id=?`

	res, err := Db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"

	res, err := Db.Exec(query, next, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Запись не найдена
	}

	return nil
}

func SearchByString(search string) ([]*Task, error) {
	var tasks []*Task = make([]*Task, 0)
	//search = "%" + strings.ToLower(search) + "%"

	query := `SELECT id, date, title, comment, repeat FROM scheduler
        WHERE title LIKE ? COLLATE NOCASE OR comment LIKE ? COLLATE NOCASE
        ORDER BY date DESC`

	rows, err := Db.Query(query, "%"+search+"%", "%"+search+"%")
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return []*Task{}, err
	}
	return tasks, nil
}

func SearchByDate(date string) ([]*Task, error) {
	var tasks []*Task = make([]*Task, 0)
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ?"

	rows, err := Db.Query(query, date)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return []*Task{}, err
	}
	return tasks, nil
}
