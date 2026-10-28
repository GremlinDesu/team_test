package usecase

import "team-test/internal/domain"

type DBPort interface {
	CreateTeam(user_id int64, team domain.Team) (int64, error)
	CreateUser(user domain.User) (int64, error)
	CreateTask(UserID int64, Task domain.Task) (int64, error)
	UpdateTask(userID int64, taskID int64, Task domain.Task) (int64, error)
	GetUserTeams(userID int64) ([]domain.Team, error)
	GetTaskHistory(userID, taskID int64) ([]domain.TaskHistory, error)
}
