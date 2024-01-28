package models

import (
	"context"
	"database/sql"
	"errors"
	"task-manger-service/db"
	"task-manger-service/domain"

	"github.com/Masterminds/squirrel"
)

func CreateTask(ctx context.Context, task domain.Task) (int, error) {

	var conn *sql.Conn
	conn, err := db.Client.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	val := []interface{}{task.Title, task.AssignedTo, task.Description, task.Priority, task.Completed, task.DueDate, task.AssignedBy}
	qbuilder := squirrel.Insert("tasks").Columns("title", "assigned_to", "description", "priority", "completed", "due_date", "assigned_by")
	qbuilder = qbuilder.Values(val...)
	query, _, err := qbuilder.ToSql()
	if err != nil {
		return 0, err
	}
	result, err := conn.ExecContext(ctx, query, val...)

	if err != nil {
		return 0, err
	}
	userId, _ := result.LastInsertId()
	return int(userId), err
}

func UpdateTask(ctx context.Context, task domain.Task) error {
	qbuilder := squirrel.Update("tasks")
	if len(task.Title) > 0 {
		qbuilder = qbuilder.Set("title", task.Title)
	}

	if task.AssignedTo > 0 {
		qbuilder = qbuilder.Set("assigned_to", task.AssignedTo)
	}
	if len(task.Description) > 0 {
		qbuilder = qbuilder.Set("description", task.Description)
	}
	if task.Priority > 0 {
		qbuilder = qbuilder.Set("priority", task.Priority)
	}

	if len(task.DueDate) > 0 {
		qbuilder = qbuilder.Set("due_date", task.DueDate)
	}

	qbuilder = qbuilder.Where(squirrel.Eq{"task_id": task.TaskId})
	query, qargs, err := qbuilder.ToSql()
	if err != nil {
		return err
	}
	conn, err := db.Client.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	result, err := conn.ExecContext(ctx, query, qargs...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if int(rowsAffected) == 0 {
		return errors.New("duplicate update  or invalid task")
	}
	return nil
}

func GetTasks(ctx context.Context, userId int, sortBy string, isAdmin bool) ([]*domain.Task, error) {
	var conn *sql.Conn
	conn, err := db.Client.Conn(ctx)

	if err != nil {
		return nil, err
	}
	defer conn.Close()
	qbuilder := squirrel.Select("task_id", "title", "assigned_to", "description", "priority", "completed", "due_date", "assigned_by").From("tasks")
	if isAdmin {
		qbuilder = qbuilder.Where(squirrel.Eq{"assigned_by": userId})
	} else {
		qbuilder = qbuilder.Where(squirrel.Eq{"assigned_to": userId})
	}

	// check If any sorting is applied

	if sortBy == "status" {
		qbuilder = qbuilder.OrderBy("completed desc")
	} else if sortBy == "dueDate" {
		qbuilder = qbuilder.OrderBy("due_date desc")
	} else if sortBy == "priority" {
		qbuilder = qbuilder.OrderBy("priority desc")
	}

	query, qargs, err := qbuilder.ToSql()
	if err != nil {
		return nil, err
	}
	userTasks := []*domain.Task{}
	rows, err := conn.QueryContext(ctx, query, qargs...)
	if err != nil {
		return userTasks, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &domain.Task{}
		err = rows.Scan(&task.TaskId, &task.Title, &task.AssignedTo, &task.Description, &task.Priority, &task.Completed, &task.DueDate, &task.AssignedBy)
		if err != nil {
			return userTasks, err
		}
		userTasks = append(userTasks, task)

	}
	return userTasks, nil

}

func MarkComplete(ctx context.Context, taskId int) error {
	qbuilder := squirrel.Update("tasks").
		Set("completed", 1).
		Where(squirrel.Eq{"task_id": taskId})
	query, qargs, err := qbuilder.ToSql()
	if err != nil {
		return err
	}
	conn, err := db.Client.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	result, err := conn.ExecContext(ctx, query, qargs...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if int(rowsAffected) == 0 {
		return errors.New("duplicate update  or invalid task")
	}
	return nil
}
