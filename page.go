package notion

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Page is a resource on the Notion platform. Its parent is either a workspace,
// another page, or a database.
// See: https://developers.notion.com/reference/page
type Page struct {
	ID             string    `json:"id"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	Parent         Parent    `json:"parent"`
	Archived       bool      `json:"archived"`
	URL            string    `json:"url"`

	// Properties differ between parent type.
	// See the `UnmarshalJSON` method.
	Properties interface{} `json:"properties"`
}

// PageProperties are properties of a page whose parent is a page or a workspace.
type PageProperties struct {
	Title PageTitle `json:"title"`
}

type PageTitle struct {
	Title []RichText `json:"title"`
}

// DatabasePageProperties are properties of a page whose parent is a database.
type DatabasePageProperties map[string]DatabasePageProperty

type DatabasePageProperty struct {
	ID   string               `json:"id,omitempty"`
	Type DatabasePropertyType `json:"type,omitempty"`

	Title          []RichText      `json:"title,omitempty"`
	RichText       []RichText      `json:"rich_text,omitempty"`
	Number         *float64        `json:"number,omitempty"`
	Select         *SelectOptions  `json:"select,omitempty"`
	MultiSelect    []SelectOptions `json:"multi_select,omitempty"`
	Date           *Date           `json:"date,omitempty"`
	Formula        *FormulaResult  `json:"formula,omitempty"`
	Relation       []Relation      `json:"relation,omitempty"`
	Rollup         *RollupResult   `json:"rollup,omitempty"`
	People         []User          `json:"people,omitempty"`
	Files          []File          `json:"files,omitempty"`
	Checkbox       *bool           `json:"checkbox,omitempty"`
	URL            *string         `json:"url,omitempty"`
	Email          *string         `json:"email,omitempty"`
	PhoneNumber    *string         `json:"phone_number,omitempty"`
	CreatedTime    *time.Time      `json:"created_time,omitempty"`
	CreatedBy      *User           `json:"created_by,omitempty"`
	LastEditedTime *time.Time      `json:"last_edited_time,omitempty"`
	LastEditedBy   *User           `json:"last_edited_by,omitempty"`
}

// CreatePageParams are the params used for creating a page.
type CreatePageParams struct {
	ParentType ParentType
	ParentID   string

	// Either DatabasePageProperties or Title must be not nil.
	DatabasePageProperties *DatabasePageProperties
	Title                  []RichText

	// Optionally, children blocks are added to the page.
	Children []Block
}

type UpdatePageParams struct {
	// Either DatabasePageProperties or Title must be not nil.
	DatabasePageProperties *DatabasePageProperties
	Title                  []RichText
}

// Value returns the underlying database page property value, based on its `type` field.
// When type is unknown/unmapped or doesn't have a value, `nil` is returned.
func (prop DatabasePageProperty) Value() interface{} {
	switch prop.Type {
	case DBPropTypeTitle:
		return prop.Title
	case DBPropTypeRichText:
		return prop.RichText
	case DBPropTypeNumber:
		return prop.Number
	case DBPropTypeSelect:
		return prop.Select
	case DBPropTypeMultiSelect:
		return prop.MultiSelect
	case DBPropTypeDate:
		return prop.Date
	case DBPropTypePeople:
		return prop.People
	case DBPropTypeFiles:
		return prop.Files
	case DBPropTypeCheckbox:
		return prop.Checkbox
	case DBPropTypeURL:
		return prop.URL
	case DBPropTypeEmail:
		return prop.Email
	case DBPropTypePhoneNumber:
		return prop.PhoneNumber
	case DBPropTypeFormula:
		return prop.Formula
	case DBPropTypeRelation:
		return prop.Relation
	case DBPropTypeRollup:
		return prop.Rollup
	case DBPropTypeCreatedTime:
		return prop.CreatedTime
	case DBPropTypeCreatedBy:
		return prop.CreatedBy
	case DBPropTypeLastEditedTime:
		return prop.LastEditedTime
	case DBPropTypeLastEditedBy:
		return prop.LastEditedBy
	default:
		return nil
	}
}

func (p CreatePageParams) Validate() error {
	if p.ParentType == "" {
		return errors.New("parent type is required")
	}
	if p.ParentID == "" {
		return errors.New("parent ID is required")
	}
	if p.ParentType == ParentTypeDatabase && p.DatabasePageProperties == nil {
		return errors.New("database page properties is required when parent type is database")
	}
	if p.ParentType == ParentTypePage && p.Title == nil {
		return errors.New("title is required when parent type is page")
	}

	return nil
}

func (p CreatePageParams) MarshalJSON() ([]byte, error) {
	type CreatePageParamsDTO struct {
		Parent     Parent      `json:"parent"`
		Properties interface{} `json:"properties"`
		Children   []Block     `json:"children,omitempty"`
	}

	var parent Parent

	if p.DatabasePageProperties != nil {
		parent.DatabaseID = p.ParentID
	} else if p.Title != nil {
		parent.PageID = p.ParentID
	}

	dto := CreatePageParamsDTO{
		Parent:   parent,
		Children: p.Children,
	}

	if p.DatabasePageProperties != nil {
		dto.Properties = p.DatabasePageProperties
	} else if p.Title != nil {
		dto.Properties = PageTitle{
			Title: p.Title,
		}
	}

	return json.Marshal(dto)
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

func (p UpdatePageParams) Validate() error {
	if p.DatabasePageProperties == nil && p.Title == nil {
		return errors.New("either database page properties or title is required")
	}
	return nil
}

func (p UpdatePageParams) MarshalJSON() ([]byte, error) {
	type UpdatePageParamsDTO struct {
		Properties interface{} `json:"properties"`
	}

	var dto UpdatePageParamsDTO

	if p.DatabasePageProperties != nil {
		dto.Properties = p.DatabasePageProperties
	} else if p.Title != nil {
		dto.Properties = PageTitle{
			Title: p.Title,
		}
	}

	return json.Marshal(dto)
}
