package models

// Comment represents a Hacker News comment with its details.
// It includes fields for ID, type, text, author, creation time, parent comment ID, and a list of child comment IDs.
type Comment struct {
	ID         int    `json:"id" db:"id"`
	Type       string `json:"type" db:"type"`
	Text       string `json:"text" db:"text"`
	Author     string `json:"by" db:"author"`
	Parent     int    `json:"parent" db:"parent_id"`
	Replies    []int  `json:"kids" db:"reply_ids"`
	Created_At int64  `json:"time" db:"created_at"`
}

// IsValid checks if the Comment has valid data.
func (c *Comment) IsValid() bool {
	return c.ID > 0 && c.Type == "Comment" && c.Text != "" && c.Author != "" && c.Created_At > 0
}
