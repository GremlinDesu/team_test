package domain

import "time"

type User struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Password_Hash string    `json:"password_hash"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
}

type Team struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBY int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamMember struct {
	ID     int64  `json:"id"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type Task struct {
	ID          int64     `json:"id"`
	TeamID      int64     `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedBy   int64     `json:"created_by"`
	AssigneeID  int64     `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    time.Time `json:"closed_at"`
	Version     string    `json:"version"`
}

type TaskHistory struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	ChangedBy string    `json:"changed_by"`
	Changes   string    `json:"changes"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskComment struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	UsedID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Analytics struct {
	TeamName          string `json:"team_name"`
	TopUsers          []User `json:"top_contributors"`
	MostCommentedTask Task   `json:"most_commented_task"`
}
