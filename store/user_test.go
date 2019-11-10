package store_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/pdk/forum/model"
	"github.com/pdk/forum/store"
)

func TestCreateQueryUser(t *testing.T) {

	db, _ := store.NewConnection(":memory:")

	createTableUsersStmt := `
	create table if not exists users (
		id integer primary key autoincrement,
		joined_at timestamp not null,
		name varchar not null unique
	);
	`

	_, err := db.Exec(createTableUsersStmt)
	if err != nil {
		t.Errorf("expected to create table users, but failed: %w", err)
	}

	user := model.NewUser("pdk")

	user, err = store.CreateUser(db, user)
	if err != nil {
		t.Errorf("expected to create user record, but failed: %w", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user.ID to be 1, but got: %d", user.ID)
	}

	foundUser, err := store.GetUserByName(db, user.Name)
	if err != nil {
		t.Errorf("expected to find user %s, but failed: %w", user.Name, err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("expected user IDs to match %d, but got %d", user.ID, foundUser.ID)
	}

	// sqlite does not preserve microseconds, so truncate for comparison. also,
	// exact timezone is not preserved, so compare in UTC.
	t1 := foundUser.JoinedAt.Truncate(time.Millisecond).UTC()
	t2 := user.JoinedAt.Truncate(time.Millisecond).UTC()
	if t1 != t2 {
		t.Errorf("expected user JoinedAt to match %s, but got %s", t2, t1)
	}

	if foundUser.Name != user.Name {
		t.Errorf("expected user Names to match %s, but got %s", user.Name, foundUser.Name)
	}

	noUser, err := store.GetUserByName(db, "nobody")
	if err != sql.ErrNoRows {
		t.Errorf("expected to get sql.ErrNoRows, but got %v", err)
	}

	emptyUser := model.User{}
	if noUser != emptyUser {
		t.Errorf("expected to get empty user, but got %v", noUser)
	}

	db.Close()
}
