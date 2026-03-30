package store

import (
	"testing"

	"kanbanboard/internal/model"
)

func TestCreateUser(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)

	user := seedUser(t, db, "Alice", "alice@test.com")

	if user.ID == "" {
		t.Error("expected user ID to be set")
	}
	if user.Name != "Alice" {
		t.Errorf("name = %q, want %q", user.Name, "Alice")
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected created_at to be set")
	}
}

func TestGetUserByEmail_found(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	seedUser(t, db, "Alice", "alice@test.com")

	user, err := GetUserByEmail(db, "alice@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("name = %q, want %q", user.Name, "Alice")
	}
}

func TestGetUserByEmail_notFound(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)

	_, err := GetUserByEmail(db, "nobody@test.com")
	if err != ErrUserNotFound {
		t.Errorf("err = %v, want ErrUserNotFound", err)
	}
}

func TestGetUserByID(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	created := seedUser(t, db, "Alice", "alice@test.com")

	user, err := GetUserByID(db, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "alice@test.com" {
		t.Errorf("email = %q, want %q", user.Email, "alice@test.com")
	}
}

func TestCountUsers(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)

	count, err := CountUsers(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}

	seedUser(t, db, "Alice", "alice@test.com")
	seedUser(t, db, "Bob", "bob@test.com")

	count, err = CountUsers(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestUpdatePassword(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	newHash := "$2a$10$newhashnewhashnewhashnewhash"
	if err := UpdatePassword(db, user.ID, newHash); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.PasswordHash != newHash {
		t.Error("password hash was not updated")
	}
}

func TestSoftDeleteUser(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	if err := SoftDeleteUser(db, user.ID); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	// User should still be fetchable by ID (for historical references)
	deleted, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Fatalf("get by ID after delete: %v", err)
	}
	if deleted.DeletedAt == nil {
		t.Error("expected deleted_at to be set")
	}
	if deleted.IsActive {
		t.Error("expected is_active to be false")
	}
	if deleted.PasswordHash != "" {
		t.Error("expected password hash to be cleared")
	}
}

func TestDeletedUserCannotLogin(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	if err := SoftDeleteUser(db, user.ID); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	// GetUserByEmail should not find deleted users
	_, err := GetUserByEmail(db, "alice@test.com")
	if err != ErrUserNotFound {
		t.Errorf("err = %v, want ErrUserNotFound (deleted user should not be findable by email)", err)
	}
}

func TestListActiveUsersBasic_excludesDeleted(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	seedUser(t, db, "Bob", "bob@test.com")

	if err := SoftDeleteUser(db, alice.ID); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	users, err := ListActiveUsersBasic(db)
	if err != nil {
		t.Fatalf("list active: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("got %d users, want 1", len(users))
	}
	if users[0].Name != "Bob" {
		t.Errorf("name = %q, want %q", users[0].Name, "Bob")
	}
}

func TestUnassignTasksForUser(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	alice := seedUser(t, db, "Alice", "alice@test.com")
	bob := seedUser(t, db, "Bob", "bob@test.com")
	project := seedProject(t, db, "Board", &alice.ID, nil)
	columns, _ := GetColumnsForProject(db, project.ID)

	// Create a task assigned to Bob
	task, err := CreateTask(db, model.Task{
		ProjectID:  project.ID,
		ColumnID:   columns[0].ID,
		CreatorID:  alice.ID,
		AssigneeID: &bob.ID,
		Title:      "Bob's Task",
		Priority:   "none",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if err := UnassignTasksForUser(db, bob.ID); err != nil {
		t.Fatalf("unassign: %v", err)
	}

	updated, err := GetTask(db, task.ID)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if updated.AssigneeID != nil {
		t.Error("expected assignee to be nil after unassign")
	}
}
