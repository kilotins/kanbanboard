package store

import (
	"testing"
	"time"
)

func TestCreateSession_returnsToken(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	token, err := CreateSession(db, user.ID, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) != 64 {
		t.Errorf("token length = %d, want 64 hex chars", len(token))
	}
}

func TestGetSession_valid(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	token, err := CreateSession(db, user.ID, 24*time.Hour)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	session, err := GetSession(db, token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.UserID != user.ID {
		t.Errorf("user ID = %q, want %q", session.UserID, user.ID)
	}
}

func TestGetSession_notFound(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)

	_, err := GetSession(db, "nonexistent-token")
	if err != ErrSessionNotFound {
		t.Errorf("err = %v, want ErrSessionNotFound", err)
	}
}

func TestGetSession_expired(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	// Create session that expires immediately
	token, err := CreateSession(db, user.ID, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = GetSession(db, token)
	if err != ErrSessionExpired {
		t.Errorf("err = %v, want ErrSessionExpired", err)
	}
}

func TestDeleteSession(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")

	token, err := CreateSession(db, user.ID, 24*time.Hour)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	if err := DeleteSession(db, token); err != nil {
		t.Fatalf("delete session: %v", err)
	}

	_, err = GetSession(db, token)
	if err != ErrSessionNotFound {
		t.Errorf("err = %v, want ErrSessionNotFound after delete", err)
	}
}
