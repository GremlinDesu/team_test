package repo_db

import "errors"

var ErrUserDoesNotExists = errors.New("user does not exist")
var ErrNotVisible = errors.New("task or team not visible for this user")
var ErrBadRole = errors.New("action not permitted for this role")
var ErrUSerNotInTeam = errors.New("user is not a member of the team")
var ErrUserAlreadyInTeam = errors.New("user already in team")
