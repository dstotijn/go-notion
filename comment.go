package notion

import (
	"encoding/json"
	"errors"
	"time"
)

// Comment represents a comment on a Notion page or block.
// See: https://developers.notion.com/reference/comment-object
type Comment struct {
	ID             string     `json:"id"`
	Parent         Parent     `json:"parent"`
	DiscussionID   string     `json:"discussion_id"`
	RichText       []RichText `json:"rich_text"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	CreatedBy      BaseUser   `json:"created_by"`
}

// CreateCommentParams are the params used for creating a comment.
type CreateCommentParams struct {
	// Either ParentPageID or DiscussionID must be non-empty. Also cannot be set
	// both at the same time.
	ParentPageID string
	DiscussionID string

	RichText []RichText
}

func (p CreateCommentParams) Validate() error {
	if p.ParentPageID == "" && p.DiscussionID == "" {
		return errors.New("either parent page ID or discussion ID is required")
	}
	if p.ParentPageID != "" && p.DiscussionID != "" {
		return errors.New("parent page ID and discussion ID cannot both be non-empty")
	}
	if len(p.RichText) == 0 {
		return errors.New("rich text is required")
	}

	return nil
}

func (p CreateCommentParams) MarshalJSON() ([]byte, error) {
	type CreateCommentParamsDTO struct {
		Parent       *Parent    `json:"parent,omitempty"`
		DiscussionID string     `json:"discussion_id,omitempty"`
		RichText     []RichText `json:"rich_text"`
	}

	dto := CreateCommentParamsDTO{
		RichText: p.RichText,
	}
	if p.ParentPageID != "" {
		dto.Parent = &Parent{
			Type:   ParentTypePage,
			PageID: p.ParentPageID,
		}
	} else {
		dto.DiscussionID = p.DiscussionID
	}

	return json.Marshal(dto)
}

// FindCommentsByBlockIDQuery is used when listing comments.
type FindCommentsByBlockIDQuery struct {
	BlockID     string
	StartCursor string
	PageSize    int
}

// FindCommentsResponse contains results (comments) and pagination data returned
// from a list request.
type FindCommentsResponse struct {
	Results    []Comment `json:"results"`
	HasMore    bool      `json:"has_more"`
	NextCursor *string   `json:"next_cursor"`
}
