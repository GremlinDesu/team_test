package repo_db

import (
	"database/sql"
	"fmt"
	"team-test/internal/domain"
)

type Repo_DB struct {
	Pool *sql.DB
}

func NewRepoDB(pool *sql.DB) *Repo_DB {
	return &Repo_DB{
		Pool: pool,
	}
}

func (r *Repo_DB) CreateUser(user domain.User) (int64, error) {
	res, err := r.Pool.Exec("INSERT INTO users (email, password_hash, name) VALUES (?, ?, ?)", user.Email, user.Password_Hash, user.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo_DB) CreateTeam(user_id int64, team domain.Team) (int64, error) {
	res, err := r.Pool.Exec("INSERT INTO teams (name, created_by) VALUES (?, ?)", team.Name, team.CreatedBY)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	makeOwnerErr := r.makeOwner(user_id, id)
	if makeOwnerErr != nil {
		return 0, makeOwnerErr
	}
	return id, nil
}

func (r *Repo_DB) CreateTask(UserID int64, Task domain.Task) (int64, error) {
	res, err := r.secureQuery("INSERT INTO tasks (team_id, title, description, status, created_by, assignee_id, version) VALUES (?, ?, ?, ?, ?, ?, ?)", Task.TeamID, Task.Title, Task.Description, Task.Status, Task.CreatedBy, Task.AssigneeID, Task.Version)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	id, err = r.updateHistory(id, UserID, "Task created")
	if err != nil {
		return 0, err
	}
	err = r.makeOwner(UserID, Task.TeamID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo_DB) UpdateTask(userID int64, taskID int64, Task domain.Task) (int64, error) {
	visible, err := r.checkVisibility(userID, taskID)
	if err != nil {
		return 0, err
	}
	if !visible {
		return 0, ErrNotVisible
	}
	return r.updateTask(taskID, Task)
}

func (r *Repo_DB) updateTask(taskID int64, Task domain.Task) (int64, error) {
	res, err := r.secureQuery("UPDATE tasks SET title = ?, description = ?, status = ?, assignee_id = ?, version = ? WHERE id = ?", Task.Title, Task.Description, Task.Status, Task.AssigneeID, Task.Version, taskID)
	if err != nil {
		return 0, err
	}
	_, err = r.updateHistory(taskID, Task.CreatedBy, "Task updated")
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repo_DB) addTeamMember(userID int64, teamUserID int64, role string) error {
	var exists bool
	r.Pool.QueryRow("SELECT EXISTS(SELECT 1 FROM team_members WHERE id = ? AND  team_id = ? )", userID, teamUserID).Scan(&exists)
	if exists {
		return ErrUserAlreadyInTeam
	}
	_, err := r.secureQuery("INSERT INTO team_members (user_id, team_id, role) VALUES (?, ?, ?)", userID, teamUserID, role)
	return err
}

func (r *Repo_DB) GetUserTeams(userID int64) ([]domain.Team, error) {
	rows, err := r.Pool.Query("SELECT id, name, created_by, created_at FROM teams WHERE id IN (SELECT team_id FROM team_members WHERE user_id = ?)", userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.CreatedBY, &team.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *Repo_DB) GetTaskHistory(userID, taskID int64) ([]domain.TaskHistory, error) {
	visible, err := r.checkVisibility(userID, taskID)
	if err != nil {
		return nil, err
	}
	if !visible {
		return nil, ErrNotVisible
	}
	rows, err := r.Pool.Query("SELECT id, task_id, changed_by, changes, created_at FROM task_history WHERE task_id = ?", taskID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var history []domain.TaskHistory
	for rows.Next() {
		var h domain.TaskHistory
		if err := rows.Scan(&h.ID, &h.TaskID, &h.ChangedBy, &h.Changes, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *Repo_DB) secureQuery(query string, args ...interface{}) (sql.Result, error) {
	tx, err := r.Pool.Begin()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo_DB) GetUserByEmail(email string) (*domain.User, error) {
	row := r.Pool.QueryRow("SELECT id, email, password_hash, name, created_at FROM users WHERE email = ?", email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Password_Hash, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo_DB) CreateTaskComment(taskID int64, userID int64, content string) error {
	visible, err := r.checkVisibility(userID, taskID)
	if err != nil {
		return err
	}
	if !visible {
		return ErrNotVisible
	}
	_, err = r.secureQuery("INSERT INTO task_comments (task_id, user_id, content) VALUES (?, ?, ?)", taskID, userID, content)
	return err
}

func (r *Repo_DB) GetTaskComments(userID, taskID int64) ([]domain.TaskComment, error) {
	visible, err := r.checkVisibility(userID, taskID)
	if err != nil {
		return nil, err
	}
	if !visible {
		return nil, ErrNotVisible
	}
	rows, err := r.Pool.Query("SELECT id, task_id, user_id, content, created_at FROM task_comments WHERE task_id = ?", taskID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var comments []domain.TaskComment
	for rows.Next() {
		var comment domain.TaskComment
		if err := rows.Scan(&comment.ID, &comment.TaskID, &comment.UsedID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *Repo_DB) updateHistory(taskID int64, changedBy int64, changes string) (int64, error) {
	res, err := r.secureQuery("INSERT INTO task_history (task_id, changed_by, changes) VALUES (?, ?, ?)", taskID, changedBy, changes)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *Repo_DB) AddUserToTeam(adderID int64, userID int64, teamID int64, role string) error {
	userExists, err := r.checkUserExists(userID)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}
	useRole, err := r.getRoleForUserInTeam(adderID, teamID)
	if err != nil {
		return err
	}
	if useRole != domain.RoleAdmin && useRole != domain.RoleOwner {
		return ErrBadRole
	}
	if role == domain.RoleAdmin && useRole != domain.RoleOwner {
		return ErrBadRole
	}
	return r.addTeamMember(userID, teamID, role)
}

func (r *Repo_DB) getRoleForUserInTeam(userID, teamID int64) (string, error) {
	row := r.Pool.QueryRow("SELECT role FROM team_members WHERE user_id = ? AND team_id = ?", userID, teamID)
	var role string
	err := row.Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUSerNotInTeam
		}
		return "", err
	}
	return role, nil
}

func (r *Repo_DB) getTeamID(taskID int64) (int64, error) {
	var id int64
	row := r.Pool.QueryRow("SELECT team_id FROM tasks WHERE id = ?", taskID)
	err := row.Scan(id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo_DB) makeOwner(userID, teamID int64) error {
	_, err := r.secureQuery("INSERT INTO team_members (user_id, team_id, role) VALUES (?, ?, ?)", userID, teamID, domain.RoleOwner)
	return err
}

func (r *Repo_DB) checkVisibility(userID int64, taskID int64) (bool, error) {
	team_id, err := r.getTeamID(taskID)
	if err != nil {
		return false, err
	}
	role, err := r.getRoleForUserInTeam(userID, team_id)
	if role == domain.RoleAdmin || role == domain.RoleOwner {
		return true, nil
	}
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ? AND  created_by = ? )"
	err = r.Pool.QueryRow(query, taskID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repo_DB) checkUserExists(userID int64) (bool, error) {
	row := r.Pool.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
