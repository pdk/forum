package store

import (
	"database/sql"
	"fmt"

	"github.com/pdk/forum/model"
)

// CreateTopic will insert a Topic into the database and return a modified Topic (ie with a new ID).
func CreateTopic(db *sql.DB, topic model.Topic) (model.Topic, error) {

	result, err := db.Exec(`insert into topics (created_by_id, name) values (?,?)`, topic.CreatedByID, topic.Name)
	if err != nil {
		return topic, fmt.Errorf("failed to save topic %s: %w", topic.Name, err)
	}

	topic.ID, err = result.LastInsertId()
	if err != nil {
		return topic, fmt.Errorf("failed to get new ID for topic %s: %w", topic.Name, err)
	}

	return topic, nil
}

// QueryTopics query all the topics and return them.
func QueryTopics(db *sql.DB) ([]model.Topic, error) {

	topicList := []model.Topic{}

	rows, err := db.Query(`select id, created_by_id, name from topics order by upper(name)`)
	if err != nil {
		return topicList, fmt.Errorf("failed to query topics: %w", err)
	}

	for rows.Next() {
		nextTopic := model.Topic{}
		err := rows.Scan(&nextTopic.ID, &nextTopic.CreatedByID, &nextTopic.Name)
		if err != nil {
			return topicList, fmt.Errorf("failed to scan a topic: %w", err)
		}

		topicList = append(topicList, nextTopic)
	}

	return topicList, nil
}

// GetTopicByID gets one topic or returns sql.ErrNoRows
func GetTopicByID(db *sql.DB, topicID int64) (model.Topic, error) {

	topic := model.Topic{}
	err := db.QueryRow(`select id, created_by_id, name from topics where id = ?`, topicID).
		Scan(&topic.ID, &topic.CreatedByID, &topic.Name)

	if err != nil {
		return topic, fmt.Errorf("cannot get topic %d: %w", topicID, err)
	}

	return topic, nil
}
