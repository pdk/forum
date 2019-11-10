package model

import "time"

// User is a human who uses this service.
type User struct {
	ID       int64
	JoinedAt time.Time
	Name     string
}

// NewUser returns a new User.
func NewUser(name string) User {
	return User{
		JoinedAt: time.Now(),
		Name:     name,
	}
}
