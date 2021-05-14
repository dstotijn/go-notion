package notion

import (
	"encoding/json"
	"fmt"
	"time"
)

// Page is a resource on the Notion platform. Its parent is either a workspace,
// another page, or a database.
// See: https://developers.notion.com/reference/page
type Page struct {
	ID             string     `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	Parent         PageParent `json:"parent"`
	Archived       bool       `json:"archived"`

	// Properties differ between parent type.
	// See the `UnmarshalJSON` method.
	Properties interface{} `json:"properties"`
}

type PageParent struct {
	Type string `json:"type"`

	PageID     *string `json:"page_id"`
	DatabaseID *string `json:"database_id"`
}

// PageProperties are properties of a page whose parent is a page or a workspace.
type PageProperties struct {
	Title struct {
		Title []RichText `json:"title"`
	} `json:"title"`
}

// DatabasePageProperties are properties of a page whose parent is a database.
type DatabasePageProperties map[string]DatabasePageProperty

type DatabasePageProperty struct {
	DatabaseProperty
	RichText    []RichText      `json:"rich_text"`
	Select      *SelectOptions  `json:"select"`
	MultiSelect []SelectOptions `json:"multi_select"`
}

// UnmarshalJSON implements json.Unmarshaler.
//
// Pages get a different Properties type based on the parent of the page.
// If parent type is `workspace` or `page_id`, PageProperties is used. Else if
// parent type is `database_id`, DatabasePageProperties is used.
func (p *Page) UnmarshalJSON(b []byte) error {
	type (
		PageAlias Page
		PageDTO   struct {
			PageAlias
			Properties json.RawMessage `json:"properties"`
		}
	)

	var dto PageDTO

	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	page := dto.PageAlias

	switch dto.Parent.Type {
	case "workspace":
		fallthrough
	case "page_id":
		var props PageProperties
		err := json.Unmarshal(dto.Properties, &props)
		if err != nil {
			return err
		}
		page.Properties = props
	case "database_id":
		var props DatabasePageProperties
		err := json.Unmarshal(dto.Properties, &props)
		if err != nil {
			return err
		}
		page.Properties = props
	default:
		return fmt.Errorf("unknown page parent type %q", dto.Parent.Type)
	}

	*p = Page(page)

	return nil
}
