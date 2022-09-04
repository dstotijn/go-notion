package notion

import "time"

// Comment is a resource on the Notion platform. It's the commit history for a page or block
// See: https://developers.notion.com/reference/comment-object
type Comment struct {
	ID             string     `json:"id"`
	DiscussionID   string     `json:"discussion_id"`
	CreatedTime    time.Time  `json:"created_time"`
	CreatedBy      *BaseUser  `json:"created_by,omitempty"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	LastEditedBy   *BaseUser  `json:"last_edited_by,omitempty"`
	Parent         Parent     `json:"parent"`
	RickText       []RichText `json:"rick_text"`
}

// CommentResponse contains
type CommentResponse struct {
	Results    []Comment `json:"results"`
	HasMore    bool      `json:"has_more"`
	NextCursor string    `json:"next_cursor"`
	Type       string    `json:"type"`
}
