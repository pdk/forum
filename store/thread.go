package store

import (
	"database/sql"
	"fmt"

	"github.com/pdk/forum/model"
)

// CreateThread will insert a Thread into the database and return a modified Thread (ie with a new ID).
func CreateThread(db *sql.DB, thread model.Thread) (model.Thread, error) {

	result, err := db.Exec(`insert into threads (topic_id, created_by_id, subject) values (?,?,?)`,
		thread.TopicID, thread.CreatedByID, thread.Subject)
	if err != nil {
		return thread, fmt.Errorf("failed to save thread %s: %w", thread.Subject, err)
	}

	thread.ID, err = result.LastInsertId()
	if err != nil {
		return thread, fmt.Errorf("failed to get new ID for thread %s: %w", thread.Subject, err)
	}

	return thread, nil
}

// QueryThreadsByTopicID returns the list of threads for a topic.
func QueryThreadsByTopicID(db *sql.DB, topicID int64) ([]model.Thread, error) {

	threadList := []model.Thread{}

	rows, err := db.Query(`select id, topic_id, created_by_id, subject from threads where topic_id = ? order by id desc`, topicID)
	if err != nil {
		return threadList, fmt.Errorf("failed to query threads: %w", err)
	}

	for rows.Next() {
		nextThread := model.Thread{}
		err := rows.Scan(&nextThread.ID, &nextThread.TopicID, &nextThread.CreatedByID, &nextThread.Subject)
		if err != nil {
			return threadList, fmt.Errorf("failed to scan a thread: %w", err)
		}

		threadList = append(threadList, nextThread)
	}

	return threadList, nil

}

// GetThreadByID gets one thread or returns sql.ErrNoRows
func GetThreadByID(db *sql.DB, threadID int64) (model.Thread, error) {

	thread := model.Thread{}
	err := db.QueryRow(`select id, topic_id, created_by_id, subject from threads where id = ?`, threadID).
		Scan(&thread.ID, &thread.TopicID, &thread.CreatedByID, &thread.Subject)

	if err != nil {
		return thread, fmt.Errorf("cannot get thread %d: %w", threadID, err)
	}

	return thread, nil
}
