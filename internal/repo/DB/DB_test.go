package repo_db_test

import (
	"team-test/config"
	"team-test/internal/domain"
	mysql_pool "team-test/internal/infrastructure/mysql"
	repo_db "team-test/internal/repo/DB"
	"testing"
)

func Setup() *repo_db.Repo_DB {
	var NewCfg = config.Config{MySQLdsn: "user:password@tcp(127.0.0.1:3306)/teams_db?parseTime=true"}
	pool, err := mysql_pool.NewPool(NewCfg)
	if err != nil {
		panic(err)
	}
	return repo_db.NewRepoDB(pool)
}

func TestCreateUser(t *testing.T) {
	RepoDB := Setup()
	user := domain.User{
		Email:         "test@example.com",
		Password_Hash: "hashed_password",
		Name:          "Test User",
	}
	err := RepoDB.CreateUser(user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
}

func TestCreateTeam(t *testing.T) {
	RepoDB := Setup()
	team := domain.Team{
		Name:      "Test Team",
		CreatedBY: "1",
	}
	err := RepoDB.CreateTeam(team)
	if err != nil {
		t.Errorf("Error creating team: %v", err)
	}
}

func TestCreateTask(t *testing.T) {
	RepoDB := Setup()
	task := domain.Task{
		TeamID:      "1",
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		CreatedBy:   "1",
		AssigneeID:  "1",
		Version:     "1.0",
	}
	err := RepoDB.CreateTask(task)
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}
}
