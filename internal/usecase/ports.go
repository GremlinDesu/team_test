package usecase

import "team-test/internal/domain"

type DBPort interface {
	CreateTeam(user_id int64, team domain.Team) (int64, error)
	CreateUser(user domain.User) (int64, error)
	CreateTask(UserID int64, Task domain.Task) (int64, error)
	UpdateTask(userID int64, taskID int64, Task domain.Task) (int64, error)
	GetUserTeams(userID int64) ([]domain.Team, error)
	GetTaskHistory(userID, taskID int64) ([]domain.TaskHistory, error)
	AddUserToTeam(adderID int64, userID int64, teamID int64, role string) error
	AddComment(taskID int64, userID int64, content string) (int64, error)
	GetAnalytics(teamID int64) (domain.Analytics, error)
}

type CachePort interface {
	LoadTeamTasks(TeamID int64, Tasks []domain.Task) error
	GetTeamTasks(TeamID int64) ([]domain.Task, error)
}

type AuthPort interface {
	CreateJwt(id int64) (string, string, error)
	ValidateJwt(tokenstr string) (int64, error)
	Refresh(Refresh string) (string, error)
	ComparePassword(password string, hash string) bool
	HashPassword(password string) (string, error)
}

type Usecase interface {
}
