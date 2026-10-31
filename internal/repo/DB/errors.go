package repo_db

import "errors"

var ErrUserDoesNotExists = errors.New("user does not exist desu")
var ErrNotVisible = errors.New("task or team not visible for this user desu")
var ErrBadRole = errors.New("action not permitted for this role desu")
var ErrUSerNotInTeam = errors.New("user is not a member of the team desu")
var ErrUserAlreadyInTeam = errors.New("user already in team desu")
