package store

import (
	"testing"
)

func TestCreateTeam(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	owner := seedUser(t, db, "Alice", "alice@test.com")

	team := seedTeam(t, db, "Team Alpha", owner.ID)

	if team.ID == "" {
		t.Error("expected team ID to be set")
	}
	if team.Name != "Team Alpha" {
		t.Errorf("name = %q, want %q", team.Name, "Team Alpha")
	}
	if team.OwnerID != owner.ID {
		t.Error("owner ID mismatch")
	}
}

func TestListTeamsForUser(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")

	seedTeam(t, db, "Alice Team", alice.ID)
	seedTeam(t, db, "Bob Team", bob.ID)

	teams, err := ListTeamsForUser(db, alice.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("got %d teams, want 1", len(teams))
	}
	if teams[0].Name != "Alice Team" {
		t.Errorf("name = %q, want %q", teams[0].Name, "Alice Team")
	}
}

func TestIsTeamMember_true(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	owner := seedUser(t, db, "Alice", "alice@test.com")
	member := seedUser(t, db, "Bob", "bob@test.com")
	team := seedTeam(t, db, "Team Alpha", owner.ID)

	if err := AddTeamMember(db, team.ID, member.ID); err != nil {
		t.Fatalf("add member: %v", err)
	}

	isMember, err := IsTeamMember(db, team.ID, member.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isMember {
		t.Error("expected Bob to be a team member")
	}
}

func TestIsTeamMember_false(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	owner := seedUser(t, db, "Alice", "alice@test.com")
	outsider := seedUser(t, db, "Charlie", "charlie@test.com")
	team := seedTeam(t, db, "Team Alpha", owner.ID)

	isMember, err := IsTeamMember(db, team.ID, outsider.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isMember {
		t.Error("expected Charlie not to be a team member")
	}
}

func TestAddTeamMember_idempotent(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	owner := seedUser(t, db, "Alice", "alice@test.com")
	member := seedUser(t, db, "Bob", "bob@test.com")
	team := seedTeam(t, db, "Team Alpha", owner.ID)

	if err := AddTeamMember(db, team.ID, member.ID); err != nil {
		t.Fatalf("first add: %v", err)
	}
	// Adding again should not error (ON CONFLICT DO NOTHING)
	if err := AddTeamMember(db, team.ID, member.ID); err != nil {
		t.Fatalf("second add (should be idempotent): %v", err)
	}

	members, err := ListTeamMembers(db, team.ID)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("got %d members, want 1", len(members))
	}
}

func TestRemoveTeamMember(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	owner := seedUser(t, db, "Alice", "alice@test.com")
	member := seedUser(t, db, "Bob", "bob@test.com")
	team := seedTeam(t, db, "Team Alpha", owner.ID)

	if err := AddTeamMember(db, team.ID, member.ID); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := RemoveTeamMember(db, team.ID, member.ID); err != nil {
		t.Fatalf("remove member: %v", err)
	}

	isMember, err := IsTeamMember(db, team.ID, member.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isMember {
		t.Error("expected Bob to no longer be a member")
	}
}

func TestCountProjectsForTeam(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	owner := seedUser(t, db, "Alice", "alice@test.com")
	team := seedTeam(t, db, "Team Alpha", owner.ID)

	count, err := CountProjectsForTeam(db, team.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}

	seedProject(t, db, "Project 1", nil, &team.ID)

	count, err = CountProjectsForTeam(db, team.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}
