package model

import (
	"time"
)

// Post is a single post by a user
type Post struct {
	ID         int64
	ThreadID   int64
	PostedByID int64
	PostedAt   time.Time
	Body       string
}

// NewPost initializes a new Post
func NewPost(threadID, postedByID int64, body string) Post {
	return Post{
		ThreadID:   threadID,
		PostedByID: postedByID,
		PostedAt:   time.Now(),
		Body:       body,
	}
}
