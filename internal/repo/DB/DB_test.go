package repo_db_test

import (
	"database/sql"
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
	err = clearDb(pool)
	if err != nil {
		panic(err)
	}
	return repo_db.NewRepoDB(pool)
}

func TestCreateUser(t *testing.T) {
	RepoDB := Setup()
	user := domain.User{
		Email:         "test12@example.com",
		Password_Hash: "hashed_password",
		Name:          "Test User",
	}
	id, err := RepoDB.CreateUser(user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	t.Logf("Created User with ID %d", id)
}

func TestCreateTeam(t *testing.T) {
	RepoDB := Setup()
	team := domain.Team{
		Name:      "Test Team",
		CreatedBY: 1,
	}
	_, err := RepoDB.CreateTeam(1, team)
	if err != nil {
		t.Errorf("Error creating team: %v", err)
	}
}

func TestCreateTask(t *testing.T) {
	RepoDB := Setup()
	task := domain.Task{
		TeamID:      1,
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		CreatedBy:   1,
		AssigneeID:  1,
		Version:     "1.0",
	}
	_, err := RepoDB.CreateTask(1, task)
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}
}

func TestGetUserTeams(t *testing.T) {
	RepoDB := Setup()
	userID := int64(1)
	team := domain.Team{
		Name:      "Test Team",
		CreatedBY: 1,
	}
	_, err := RepoDB.CreateTeam(1, team)
	if err != nil {
		t.Errorf("Error creating team: %v", err)
	}
	team1 := domain.Team{
		Name:      "Test11 Team",
		CreatedBY: 1,
	}
	_, err = RepoDB.CreateTeam(1, team1)
	if err != nil {
		t.Errorf("Error creating team: %v", err)
	}
	teams, err := RepoDB.GetUserTeams(userID)
	if err != nil {
		t.Errorf("Error getting user teams: %v", err)
	}
	if len(teams) == 0 {
		t.Errorf("Expected at least one team for user %d, got 0", userID)
	}
}

func TestAddMemberToTeam(t *testing.T) {
	RepoDB := Setup()
	user := domain.User{
		Email:         "test1445@example.com",
		Password_Hash: "hashed_password",
		Name:          "Test User",
	}
	id, err := RepoDB.CreateUser(user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	team := domain.Team{
		Name:      "Test Team",
		CreatedBY: id,
	}

	user2 := domain.User{
		Email:         "test1545@example.com",
		Password_Hash: "hashed_password",
		Name:          "Test User",
	}
	id2, err := RepoDB.CreateUser(user2)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	teamID, err := RepoDB.CreateTeam(id, team)
	if err != nil {
		t.Errorf("Error creating team: %v", err)
	}
	err = RepoDB.AddUserToTeam(id, id2, teamID, domain.RoleMember)
	if err != nil {
		t.Errorf("Error adding member to team: %v", err)
	}
	err = RepoDB.AddUserToTeam(id2, id, teamID, domain.RoleMember)
	if err != repo_db.ErrBadRole {
		t.Errorf("bad role not found!!")
	}
	err = RepoDB.AddUserToTeam(id, id2, teamID, domain.RoleMember)
	if err != repo_db.ErrUserAlreadyInTeam {
		t.Errorf("double add avaliable %v", err)
	}
}

func clearDb(db *sql.DB) error {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE teams")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE team_members")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE tasks")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE task_history")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE task_comments")
	if err != nil {
		return err
	}
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		return err
	}
	return nil
}
