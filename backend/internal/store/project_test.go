package store

import (
	"testing"

	"kanbanboard/internal/model"
)

func TestCreateProject_userOwned(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	project := seedProject(t, db, "My Board", &user.ID, nil)

	if project.ID == "" {
		t.Error("expected project ID to be set")
	}
	if *project.OwnerUserID != user.ID {
		t.Error("owner user ID mismatch")
	}
	if project.OwnerTeamID != nil {
		t.Error("expected no team owner")
	}
}

func TestCreateProject_teamOwned(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	team := seedTeam(t, db, "Team Alpha", user.ID)

	project := seedProject(t, db, "Team Board", nil, &team.ID)

	if project.OwnerUserID != nil {
		t.Error("expected no user owner")
	}
	if *project.OwnerTeamID != team.ID {
		t.Error("owner team ID mismatch")
	}
}

func TestCreateDefaultColumns(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "My Board", &user.ID, nil)

	columns, err := GetColumnsForProject(db, project.ID)
	if err != nil {
		t.Fatalf("get columns: %v", err)
	}
	if len(columns) != 5 {
		t.Fatalf("got %d columns, want 5", len(columns))
	}

	expected := []string{"Inbox", "Todo", "In Progress", "Blocked", "Done"}
	for i, col := range columns {
		if col.Name != expected[i] {
			t.Errorf("column %d name = %q, want %q", i, col.Name, expected[i])
		}
		if col.Position != i {
			t.Errorf("column %d position = %d, want %d", i, col.Position, i)
		}
	}
}

func TestCreateDefaultLabels(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "My Board", &user.ID, nil)

	labels, err := GetLabelsForProject(db, project.ID)
	if err != nil {
		t.Fatalf("get labels: %v", err)
	}
	if len(labels) != 4 {
		t.Fatalf("got %d labels, want 4", len(labels))
	}
}

func TestListProjectsForUser_ownerSees(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")

	seedProject(t, db, "Alice Project", &alice.ID, nil)
	seedProject(t, db, "Bob Project", &bob.ID, nil)

	projects, err := ListProjectsForUser(db, alice.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("got %d projects, want 1", len(projects))
	}
	if projects[0].Name != "Alice Project" {
		t.Errorf("name = %q, want %q", projects[0].Name, "Alice Project")
	}
}

func TestListProjectsForUser_teamMemberSees(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")

	team := seedTeam(t, db, "Team Alpha", alice.ID)
	if err := AddTeamMember(db, team.ID, bob.ID); err != nil {
		t.Fatalf("add member: %v", err)
	}

	seedProject(t, db, "Team Project", nil, &team.ID)

	projects, err := ListProjectsForUser(db, bob.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("got %d projects, want 1", len(projects))
	}
	if projects[0].Name != "Team Project" {
		t.Errorf("name = %q, want %q", projects[0].Name, "Team Project")
	}
}

func TestListProjectsForUser_publicVisible(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")

	// Create a public project owned by Alice
	project, err := CreateProject(db, model.Project{
		Name:        "Public Board",
		Visibility:  "public",
		OwnerUserID: &alice.ID,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	_ = CreateDefaultColumns(db, project.ID)

	projects, err := ListProjectsForUser(db, bob.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("got %d projects, want 1", len(projects))
	}
	if projects[0].Name != "Public Board" {
		t.Errorf("name = %q, want %q", projects[0].Name, "Public Board")
	}
}

func TestListProjectsForUser_privateHidden(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")

	seedProject(t, db, "Private Board", &alice.ID, nil) // private by default

	projects, err := ListProjectsForUser(db, bob.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("got %d projects, want 0 (private should be hidden)", len(projects))
	}
}

func TestListProjectsForUser_noDuplicates(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")

	// Alice owns a team and is also a member — the project should appear once
	team := seedTeam(t, db, "Team Alpha", alice.ID)
	if err := AddTeamMember(db, team.ID, alice.ID); err != nil {
		t.Fatalf("add member: %v", err)
	}

	project, err := CreateProject(db, model.Project{
		Name:        "Team Public",
		Visibility:  "public",
		OwnerTeamID: &team.ID,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	_ = CreateDefaultColumns(db, project.ID)

	projects, err := ListProjectsForUser(db, alice.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Errorf("got %d projects, want 1 (DISTINCT should prevent duplicates)", len(projects))
	}
}

func TestGetProjectMembers_userOwned(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")

	project := seedProject(t, db, "Alice Board", &alice.ID, nil)

	members, err := GetProjectMembers(db, project)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("got %d members, want 1", len(members))
	}
	if members[0].ID != alice.ID {
		t.Error("expected owner to be the only member")
	}
}

func TestGetProjectMembers_teamOwned(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")

	team := seedTeam(t, db, "Team Alpha", alice.ID)
	if err := AddTeamMember(db, team.ID, bob.ID); err != nil {
		t.Fatalf("add member: %v", err)
	}

	project := seedProject(t, db, "Team Board", nil, &team.ID)

	members, err := GetProjectMembers(db, project)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should include team owner (Alice) + team member (Bob)
	if len(members) != 2 {
		t.Fatalf("got %d members, want 2", len(members))
	}
}

func TestReorderColumns(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "My Board", &user.ID, nil)

	columns, err := GetColumnsForProject(db, project.ID)
	if err != nil {
		t.Fatalf("get columns: %v", err)
	}

	// Reverse the column order
	reversed := make([]string, len(columns))
	for i, col := range columns {
		reversed[len(columns)-1-i] = col.ID
	}

	if err := ReorderColumns(db, project.ID, reversed); err != nil {
		t.Fatalf("reorder: %v", err)
	}

	updated, err := GetColumnsForProject(db, project.ID)
	if err != nil {
		t.Fatalf("get columns after reorder: %v", err)
	}

	// First column should now be "Done" (was last)
	if updated[0].Name != "Done" {
		t.Errorf("first column = %q, want %q", updated[0].Name, "Done")
	}
	if updated[4].Name != "Inbox" {
		t.Errorf("last column = %q, want %q", updated[4].Name, "Inbox")
	}

	// Verify positions are sequential
	for i, col := range updated {
		if col.Position != i {
			t.Errorf("column %d position = %d, want %d", i, col.Position, i)
		}
	}
}
