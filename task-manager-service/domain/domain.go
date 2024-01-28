package domain

import (
	"errors"
)

type Task struct {
	TaskId      int    `json:"taskId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	DueDate     string `json:"dueDate"`
	Completed   bool   `json:"completed"`
	AssignedTo  int    `json:"assignedTo"` // Username of the user who owns the task
	AssignedBy  int    `json:"assignedBy"`
}

func (t *Task) Validate() error {
	if len(t.Title) <= 0 {
		return errors.New("invalid title")
	}
	if t.Priority < 0 {
		return errors.New("invalid priority")
	}

	if t.AssignedTo < 0 {
		return errors.New("invalid AssignedTo")
	}

	if t.AssignedBy < 0 {
		return errors.New("invalid AssignedBy")
	}
	return nil

}

type User struct {
	UserId    int    `json:"userId,omitempty"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Type      string `json:"type"` // "default" or "admin"
	CreatedBy int    `json:"createdBy"`
}

type TaskSearchResponse struct {
	Tasks  []*Task `json:"tasks"`
	Status string  `json:"status"`
}
type TaskResponse struct {
	Status    string `json:"status"`
	RespError string `json:"error,omitempty"`
}

type TaskCreateResponse struct {
	Status string `json:"status"`
	TaskId int    `json:"TaskId"`
}
