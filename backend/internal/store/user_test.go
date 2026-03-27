package store

import (
	"testing"
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
