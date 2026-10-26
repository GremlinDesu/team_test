package domain

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Password_Hash string    `json:"password_hash"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
}

type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedBY string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamMember struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	TeamUserID string `json:"team_user_id"`
	Role       string `json:"role"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	AssigneeID  string    `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    time.Time `json:"closed_at"`
	Version     string    `json:"version"`
}

type TaskHistory struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	ChangedBy string    `json:"changed_by"`
	Changes   string    `json:"changes"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskComment struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	UsedID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
