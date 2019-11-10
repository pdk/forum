package store

import (
	"database/sql"
	"fmt"

	"github.com/pdk/forum/model"
)

// CreateUser will insert a User into the database and return a modified User (ie with a new ID).
func CreateUser(db *sql.DB, user model.User) (model.User, error) {

	result, err := db.Exec(`insert into users (joined_at, name) values (?,?)`, user.JoinedAt, user.Name)
	if err != nil {
		return user, fmt.Errorf("failed to save user %s: %w", user.Name, err)
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return user, fmt.Errorf("failed to get new ID for user %s: %w", user.Name, err)
	}

	return user, nil
}

// GetUserByID will query and return a User by ID. If no user matches,
// sql.ErrNoRows will be returned as the error.
func GetUserByID(db *sql.DB, userID int64) (model.User, error) {

	user := model.User{}

	err := db.QueryRow(`select id, joined_at, name from users where id = ?`, userID).
		Scan(&user.ID, &user.JoinedAt, &user.Name)

	return user, err
}

// GetUserByName will query and return a User by ID. If no user matches,
// sql.ErrNoRows will be returned as the error.
func GetUserByName(db *sql.DB, name string) (model.User, error) {

	user := model.User{}

	err := db.QueryRow(`select id, joined_at, name from users where name = ?`, name).
		Scan(&user.ID, &user.JoinedAt, &user.Name)

	return user, err
}

// GetOrCreateUserByName will return either an existing user, or a newly created
// user, with the given name.
func GetOrCreateUserByName(db *sql.DB, name string) (model.User, error) {

	user, err := GetUserByName(db, name)
	if err == nil {
		return user, nil
	}

	if err != sql.ErrNoRows {
		return model.User{}, fmt.Errorf("failed to get/create user %s: %w", name, err)
	}

	user = model.NewUser(name)

	user, err = CreateUser(db, user)
	if err != nil {
		return user, fmt.Errorf("failed to get/create user %s: %w", name, err)
	}

	return user, nil
}
