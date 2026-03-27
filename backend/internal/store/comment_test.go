package store

import (
	"testing"

	"kanbanboard/internal/model"
)

func TestCreateComment(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "Board", &user.ID, nil)
	columns, _ := GetColumnsForProject(db, project.ID)
	task := seedTask(t, db, project.ID, columns[0].ID, user.ID, "Task")

	comment, err := CreateComment(db, model.Comment{
		TaskID:   task.ID,
		AuthorID: user.ID,
		Text:     "Hello world",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID == "" {
		t.Error("expected comment ID to be set")
	}
	if comment.Text != "Hello world" {
		t.Errorf("text = %q, want %q", comment.Text, "Hello world")
	}
}

func TestListCommentsForTask(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "Board", &user.ID, nil)
	columns, _ := GetColumnsForProject(db, project.ID)
	task := seedTask(t, db, project.ID, columns[0].ID, user.ID, "Task")

	CreateComment(db, model.Comment{TaskID: task.ID, AuthorID: user.ID, Text: "First"})
	CreateComment(db, model.Comment{TaskID: task.ID, AuthorID: user.ID, Text: "Second"})

	comments, err := ListCommentsForTask(db, task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("got %d comments, want 2", len(comments))
	}
	// Should include author name from JOIN
	if comments[0].AuthorName != "Alice" {
		t.Errorf("author name = %q, want %q", comments[0].AuthorName, "Alice")
	}
	// Should be ordered by created_at
	if comments[0].Text != "First" {
		t.Errorf("first comment = %q, want %q", comments[0].Text, "First")
	}
}

func TestGetComment_found(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "Board", &user.ID, nil)
	columns, _ := GetColumnsForProject(db, project.ID)
	task := seedTask(t, db, project.ID, columns[0].ID, user.ID, "Task")

	created, _ := CreateComment(db, model.Comment{TaskID: task.ID, AuthorID: user.ID, Text: "Test"})

	comment, err := GetComment(db, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Text != "Test" {
		t.Errorf("text = %q, want %q", comment.Text, "Test")
	}
}

func TestGetComment_notFound(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)

	_, err := GetComment(db, "00000000-0000-0000-0000-000000000000")
	if err != ErrCommentNotFound {
		t.Errorf("err = %v, want ErrCommentNotFound", err)
	}
}

func TestUpdateComment(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "Board", &user.ID, nil)
	columns, _ := GetColumnsForProject(db, project.ID)
	task := seedTask(t, db, project.ID, columns[0].ID, user.ID, "Task")

	created, _ := CreateComment(db, model.Comment{TaskID: task.ID, AuthorID: user.ID, Text: "Original"})

	updated, err := UpdateComment(db, created.ID, "Modified")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Text != "Modified" {
		t.Errorf("text = %q, want %q", updated.Text, "Modified")
	}
}

func TestDeleteComment(t *testing.T) {
	db := testDB(t)
	cleanTables(t, db)
	user := seedUser(t, db, "Alice", "alice@test.com")
	project := seedProject(t, db, "Board", &user.ID, nil)
	columns, _ := GetColumnsForProject(db, project.ID)
	task := seedTask(t, db, project.ID, columns[0].ID, user.ID, "Task")

	created, _ := CreateComment(db, model.Comment{TaskID: task.ID, AuthorID: user.ID, Text: "To Delete"})

	if err := DeleteComment(db, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := GetComment(db, created.ID)
	if err != ErrCommentNotFound {
		t.Errorf("err = %v, want ErrCommentNotFound", err)
	}
}
