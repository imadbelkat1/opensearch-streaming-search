package models

// Story represents a Hacker News story with its details.
// It includes fields for ID, type, title, URL, score, author, creation time
type Story struct {
	ID             int    `json:"id" db:"id"`
	Type           string `json:"type" db:"type"`
	Title          string `json:"title" db:"title"`
	URL            string `json:"url" db:"url"`
	Score          int    `json:"score" db:"score"`
	Author         string `json:"by" db:"author"`
	Created_At     int64  `json:"time" db:"created_at"`
	Comments_count int    `json:"kids" db:"comments_count"`
}

// IsValid checks if the Story has valid data.
func (s *Story) IsValid() bool {
	return s.ID > 0 && s.Type == "Story" && s.Title != "" && s.Author != "" && s.Created_At > 0
}
