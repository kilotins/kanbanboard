package handler

import "kanbanboard/internal/model"

// CanViewProject determines if a user can view a project based on ownership,
// team membership, and visibility. The caller resolves teamOwnerID and
// isTeamMember from the database before calling this function.
func CanViewProject(project model.Project, userID string, teamOwnerID string, isTeamMember bool) bool {
	// User owner can always view
	if project.OwnerUserID != nil && *project.OwnerUserID == userID {
		return true
	}

	// Public projects are visible to everyone
	if project.Visibility == "public" {
		return true
	}

	// Team owner or member can view team projects
	if project.OwnerTeamID != nil {
		if teamOwnerID == userID {
			return true
		}
		return isTeamMember
	}

	return false
}

// CanEditProject determines if a user can edit tasks in a project.
// Unlike view, public visibility does NOT grant edit access.
func CanEditProject(project model.Project, userID string, teamOwnerID string, isTeamMember bool) bool {
	// User owner can edit
	if project.OwnerUserID != nil && *project.OwnerUserID == userID {
		return true
	}

	// Team owner or member can edit team projects
	if project.OwnerTeamID != nil {
		if teamOwnerID == userID {
			return true
		}
		return isTeamMember
	}

	return false
}

// IsProjectOwner checks if a user is the direct owner or the team owner of a project.
func IsProjectOwner(project model.Project, userID string, teamOwnerID string) bool {
	if project.OwnerUserID != nil && *project.OwnerUserID == userID {
		return true
	}
	if project.OwnerTeamID != nil && teamOwnerID == userID {
		return true
	}
	return false
}

// IsTeamOwner checks if a user owns a team.
func IsTeamOwner(team model.Team, userID string) bool {
	return team.OwnerID == userID
}

// IsCommentAuthor checks if a user is the author of a comment.
func IsCommentAuthor(comment model.Comment, userID string) bool {
	return comment.AuthorID == userID
}
