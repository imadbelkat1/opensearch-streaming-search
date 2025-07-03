package models

// Poll represents a Hacker News poll with its details.
// It includes fields for ID, type, title, score, author, creation time, description, and a list of part IDs.
type Poll struct {
	ID            int    `json:"id" db:"id"`
	Type          string `json:"type" db:"type"`
	Title         string `json:"title" db:"title"`
	Score         int    `json:"score" db:"score"`
	Author        string `json:"by" db:"author"`
	Created_At    int64  `json:"time" db:"created_at"`
	Comment_count string `json:"desc" db:"desc"`
	Reply_Ids     []int  `json:"reply_ids" db:"reply_ids"` // Array of IDs for replies to this poll
	Parts         []int  `json:"parts" db:"poll_options"`
}

// IsValid checks if the Poll has valid data.
func (p *Poll) IsValid() bool {
	return p.ID > 0 && p.Type == "Poll" && p.Title != "" && p.Author != "" && p.Created_At > 0
}
