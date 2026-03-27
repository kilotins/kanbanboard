package store

import (
	"database/sql"
	"fmt"
)

// IsTeamMember checks if a user is a member of a team.
func IsTeamMember(db *sql.DB, teamID, userID string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)",
		teamID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check team membership: %w", err)
	}
	return exists, nil
}
