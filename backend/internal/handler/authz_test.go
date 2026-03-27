package handler

import (
	"testing"

	"kanbanboard/internal/model"
)

func strPtr(s string) *string { return &s }

// --- CanViewProject ---

func TestCanViewProject_ownerCanView(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1"), Visibility: "private"}
	if !CanViewProject(p, "u1", "", false) {
		t.Error("owner should be able to view their own project")
	}
}

func TestCanViewProject_publicVisibleToAll(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1"), Visibility: "public"}
	if !CanViewProject(p, "u2", "", false) {
		t.Error("public project should be visible to any user")
	}
}

func TestCanViewProject_privateHiddenFromNonOwner(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1"), Visibility: "private"}
	if CanViewProject(p, "u2", "", false) {
		t.Error("private project should be hidden from non-owner")
	}
}

func TestCanViewProject_teamOwnerCanView(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "private"}
	if !CanViewProject(p, "u1", "u1", false) {
		t.Error("team owner should be able to view team project")
	}
}

func TestCanViewProject_teamMemberCanView(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "private"}
	if !CanViewProject(p, "u2", "u1", true) {
		t.Error("team member should be able to view team project")
	}
}

func TestCanViewProject_nonMemberCannotViewPrivateTeamProject(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "private"}
	if CanViewProject(p, "u3", "u1", false) {
		t.Error("non-member should not view private team project")
	}
}

func TestCanViewProject_publicTeamProjectVisibleToAll(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "public"}
	if !CanViewProject(p, "u3", "u1", false) {
		t.Error("public team project should be visible to all")
	}
}

// --- CanEditProject ---

func TestCanEditProject_ownerCanEdit(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1"), Visibility: "private"}
	if !CanEditProject(p, "u1", "", false) {
		t.Error("owner should be able to edit their own project")
	}
}

func TestCanEditProject_nonOwnerCannotEdit(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1"), Visibility: "private"}
	if CanEditProject(p, "u2", "", false) {
		t.Error("non-owner should not be able to edit project")
	}
}

func TestCanEditProject_publicDoesNotGrantEdit(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1"), Visibility: "public"}
	if CanEditProject(p, "u2", "", false) {
		t.Error("public visibility should NOT grant edit access")
	}
}

func TestCanEditProject_teamOwnerCanEdit(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "private"}
	if !CanEditProject(p, "u1", "u1", false) {
		t.Error("team owner should be able to edit team project")
	}
}

func TestCanEditProject_teamMemberCanEdit(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "private"}
	if !CanEditProject(p, "u2", "u1", true) {
		t.Error("team member should be able to edit team project")
	}
}

func TestCanEditProject_nonMemberCannotEditTeamProject(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1"), Visibility: "public"}
	if CanEditProject(p, "u3", "u1", false) {
		t.Error("non-member should not be able to edit team project even if public")
	}
}

// --- IsProjectOwner ---

func TestIsProjectOwner_directOwner(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1")}
	if !IsProjectOwner(p, "u1", "") {
		t.Error("direct owner should be recognized")
	}
}

func TestIsProjectOwner_nonOwner(t *testing.T) {
	p := model.Project{OwnerUserID: strPtr("u1")}
	if IsProjectOwner(p, "u2", "") {
		t.Error("non-owner should not be recognized as owner")
	}
}

func TestIsProjectOwner_teamOwner(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1")}
	if !IsProjectOwner(p, "u1", "u1") {
		t.Error("team owner should be recognized as project owner")
	}
}

func TestIsProjectOwner_teamMemberNotOwner(t *testing.T) {
	p := model.Project{OwnerTeamID: strPtr("t1")}
	if IsProjectOwner(p, "u2", "u1") {
		t.Error("team member who is not team owner should not be project owner")
	}
}

// --- IsTeamOwner ---

func TestIsTeamOwner_owner(t *testing.T) {
	team := model.Team{OwnerID: "u1"}
	if !IsTeamOwner(team, "u1") {
		t.Error("should recognize team owner")
	}
}

func TestIsTeamOwner_nonOwner(t *testing.T) {
	team := model.Team{OwnerID: "u1"}
	if IsTeamOwner(team, "u2") {
		t.Error("should not recognize non-owner as team owner")
	}
}

// --- IsCommentAuthor ---

func TestIsCommentAuthor_author(t *testing.T) {
	c := model.Comment{AuthorID: "u1"}
	if !IsCommentAuthor(c, "u1") {
		t.Error("should recognize comment author")
	}
}

func TestIsCommentAuthor_nonAuthor(t *testing.T) {
	c := model.Comment{AuthorID: "u1"}
	if IsCommentAuthor(c, "u2") {
		t.Error("should not recognize non-author")
	}
}
