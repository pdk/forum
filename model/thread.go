package model

// Thread is a chain of posts with a single subject.
type Thread struct {
	ID          int64
	TopicID     int64
	CreatedByID int64
	Subject     string
}

// NewThread returns a new Thread.
func NewThread(topicID, creatorID int64, subject string) Thread {
	return Thread{
		TopicID:     topicID,
		CreatedByID: creatorID,
		Subject:     subject,
	}
}
