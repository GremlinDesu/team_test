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

func (r *Repo_DB) CreateUser(user domain.User) error {
	_, err := r.Pool.Exec("INSERT INTO users (email, password_hash, name) VALUES (?, ?, ?)", user.Email, user.Password_Hash, user.Name)
	return err
}

func (r *Repo_DB) CreateTeam(team domain.Team) error {
	_, err := r.Pool.Exec("INSERT INTO teams (name, created_by) VALUES (?, ?)", team.Name, team.CreatedBY)
	return err
}

func (r *Repo_DB) CreateTask(Task domain.Task) error {
	var id string
	_, err := r.secureQuery("INSERT INTO tasks (team_id, title, description, status, created_by, assignee_id, version) VALUES (?, ?, ?, ?, ?, ?, ?)", Task.TeamID, Task.Title, Task.Description, Task.Status, Task.CreatedBy, Task.AssigneeID, Task.Version)
	if err != nil {
		return err
	}
	row, err := r.Pool.Query("SELECT LAST_INSERT_ID()")
	if err != nil {
		return err
	}
	defer row.Close()
	if row.Next() {
		err = row.Scan(&id)
	}
	if err != nil {
		return err
	}
	return r.updateHistory(id, Task.CreatedBy, "Task created")
}

func (r *Repo_DB) UpdateTask(AdderID string, taskID string, Task domain.Task) error {

	return r.updateTask(taskID, Task)
}

func (r *Repo_DB) updateTask(taskID string, Task domain.Task) error {
	_, err := r.secureQuery("UPDATE tasks SET title = ?, description = ?, status = ?, assignee_id = ?, version = ? WHERE id = ?", Task.Title, Task.Description, Task.Status, Task.AssigneeID, Task.Version, taskID)
	if err != nil {
		return err
	}
	return r.updateHistory(taskID, Task.CreatedBy, "Task updated")
}

func (r *Repo_DB) addTeamMember(userID, teamUserID, role string) error {
	_, err := r.secureQuery("INSERT INTO team_members (user_id, team_user_id, role) VALUES (?, ?, ?)", userID, teamUserID, role)
	return err
}

func (r *Repo_DB) GetUserTeams(userID string) ([]domain.Team, error) {
	rows, err := r.Pool.Query("SELECT t.id, t.name, t.created_by, t.created_at FROM teams t JOIN team_members tm ON t.id = tm.team_user_id WHERE tm.user_id = ?", userID)
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

func (r *Repo_DB) GetTaskHistory(taskID int64) ([]domain.TaskHistory, error) {
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

func (r *Repo_DB) secureQuery(query string, args ...interface{}) (*sql.Rows, error) {
	tx, err := r.Pool.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if rows.Err() != nil {
		tx.Rollback()
		return nil, rows.Err()
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return rows, nil
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

func (r *Repo_DB) CreateTaskComment(taskID int64, userID, content string) error {
	_, err := r.secureQuery("INSERT INTO task_comments (task_id, user_id, content) VALUES (?, ?, ?)", taskID, userID, content)
	return err
}

func (r *Repo_DB) GetTaskComments(taskID int64) ([]domain.TaskComment, error) {
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

func (r *Repo_DB) updateHistory(taskID string, changedBy string, changes string) error {
	_, err := r.secureQuery("INSERT INTO task_history (task_id, changed_by, changes) VALUES (?, ?, ?)", taskID, changedBy, changes)
	return err
}

func (r *Repo_DB) AddUserToTeam(adderID, userID, teamID, role string) error {
	useRole, err := r.GetRoleForUserInTeam(adderID, teamID)
	if err != nil {
		return err
	}
	if useRole != domain.RoleAdmin && useRole != domain.RoleOwner {
		return fmt.Errorf("user does not have permission to add members to this team")
	}
	if role == domain.RoleAdmin && useRole != domain.RoleOwner {
		return fmt.Errorf("only owners can assign admin role")
	}
	return r.addTeamMember(userID, teamID, role)
}

func (r *Repo_DB) GetTaskByID(taskID string) (*domain.Task, error) {
	row := r.Pool.QueryRow("SELECT id, team_id, title, description, status, created_by, assignee_id, created_at, updated_at, closed_at, version FROM tasks WHERE id = ?", taskID)
	var task domain.Task
	err := row.Scan(&task.ID, &task.TeamID, &task.Title, &task.Description, &task.Status, &task.CreatedBy, &task.AssigneeID, &task.CreatedAt, &task.UpdatedAt, &task.ClosedAt, &task.Version)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repo_DB) GetRoleForUserInTeam(userID, teamID string) (string, error) {
	row := r.Pool.QueryRow("SELECT role FROM team_members WHERE user_id = ? AND team_user_id = ?", userID, teamID)
	var role string
	err := row.Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}
