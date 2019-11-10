package model

// Topic is an area of discussion.
type Topic struct {
	ID          int64
	CreatedByID int64
	Name        string
}

// NewTopic makes a new Topic.
func NewTopic(creatorID int64, name string) Topic {
	return Topic{
		CreatedByID: creatorID,
		Name:        name,
	}
}
