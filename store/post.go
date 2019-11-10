package store

import (
	"database/sql"
	"fmt"

	"github.com/pdk/forum/model"
)

// CreatePost will insert a Post into the database and return a modified Post (ie with a new ID).
func CreatePost(db *sql.DB, post model.Post) (model.Post, error) {

	result, err := db.Exec(`insert into posts (thread_id, posted_by_id, posted_at, body) values (?,?,?,?)`,
		post.ThreadID, post.PostedByID, post.PostedAt, post.Body)
	if err != nil {
		return post, fmt.Errorf("failed to save post %s: %w", post.Body, err)
	}

	post.ID, err = result.LastInsertId()
	if err != nil {
		return post, fmt.Errorf("failed to get new ID for post %s: %w", post.Body, err)
	}

	return post, nil
}

// QueryPostsByThreadID selects all the posts for a given thread.
func QueryPostsByThreadID(db *sql.DB, threadID int64) ([]model.Post, error) {

	postList := []model.Post{}

	rows, err := db.Query(`select id, thread_id, posted_by_id, posted_at, body from posts where thread_id = ? order by id asc`, threadID)
	if err != nil {
		return postList, fmt.Errorf("failed to query posts by id %d: %w", threadID, err)
	}

	for rows.Next() {
		nextPost := model.Post{}
		err := rows.Scan(&nextPost.ID, &nextPost.ThreadID, &nextPost.PostedByID, &nextPost.PostedAt, &nextPost.Body)
		if err != nil {
			return postList, fmt.Errorf("failed to scan post row: %w", err)
		}

		postList = append(postList, nextPost)
	}

	return postList, nil
}
