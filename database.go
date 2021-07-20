package notion

import (
	"encoding/json"
	"errors"
	"time"
)

// Database is a resource on the Notion platform.
// See: https://developers.notion.com/reference/database
type Database struct {
	ID             string             `json:"id"`
	CreatedTime    time.Time          `json:"created_time"`
	LastEditedTime time.Time          `json:"last_edited_time"`
	Title          []RichText         `json:"title"`
	Properties     DatabaseProperties `json:"properties"`
	Parent         Parent             `json:"parent"`
}

// DatabaseProperties is a mapping of properties defined on a database.
type DatabaseProperties map[string]DatabaseProperty

// Database property metadata types.
type (
	EmptyMetadata  struct{}
	NumberMetadata struct {
		Format NumberFormat `json:"format"`
	}
	SelectMetadata struct {
		Options []SelectOptions `json:"options"`
	}
	FormulaMetadata struct {
		Expression string `json:"expression"`
	}
	RelationMetadata struct {
		DatabaseID     string `json:"database_id,omitempty"`
		SyncedPropName string `json:"synced_property_name,omitempty"`
		SyncedPropID   string `json:"synced_property_id,omitempty"`
	}
	RollupMetadata struct {
		RelationPropName string `json:"relation_property_name,omitempty"`
		RelationPropID   string `json:"relation_property_id,omitempty"`
		RollupPropName   string `json:"rollup_property_name,omitempty"`
		RollupPropID     string `json:"rollup_property_id,omitempty"`
		Function         string `json:"function,omitempty"`
	}
)

type SelectOptions struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color Color  `json:"color,omitempty"`
}

type FormulaResult struct {
	Type FormulaResultType `json:"type"`

	String  *string  `json:"string,omitempty"`
	Number  *float64 `json:"number,omitempty"`
	Boolean *bool    `json:"boolean,omitempty"`
	Date    *Date    `json:"date,omitempty"`
}

type Relation struct {
	ID string `json:"id"`
}

type RollupResult struct {
	Type RollupResultType `json:"type"`

	Number *float64               `json:"number,omitempty"`
	Date   *Date                  `json:"date,omitempty"`
	Array  []DatabasePageProperty `json:"array,omitempty"`
}

type People struct {
	People []User `json:"people"`
}

type File struct {
	Name string `json:"name"`
}

type DatabaseProperty struct {
	ID   string               `json:"id,omitempty"`
	Type DatabasePropertyType `json:"type"`

	Title          *EmptyMetadata `json:"title,omitempty"`
	RichText       *EmptyMetadata `json:"rich_text,omitempty"`
	Date           *EmptyMetadata `json:"date,omitempty"`
	People         *EmptyMetadata `json:"people,omitempty"`
	Files          *EmptyMetadata `json:"files,omitempty"`
	Checkbox       *EmptyMetadata `json:"checkbox,omitempty"`
	URL            *EmptyMetadata `json:"url,omitempty"`
	Email          *EmptyMetadata `json:"email,omitempty"`
	PhoneNumber    *EmptyMetadata `json:"phone_number,omitempty"`
	CreatedTime    *EmptyMetadata `json:"created_time,omitempty"`
	CreatedBy      *EmptyMetadata `json:"created_by,omitempty"`
	LastEditedTime *EmptyMetadata `json:"last_edited_time,omitempty"`
	LastEditedBy   *EmptyMetadata `json:"last_edited_by,omitempty"`

	Number      *NumberMetadata   `json:"number,omitempty"`
	Select      *SelectMetadata   `json:"select,omitempty"`
	MultiSelect *SelectMetadata   `json:"multi_select,omitempty"`
	Formula     *FormulaMetadata  `json:"formula,omitempty"`
	Relation    *RelationMetadata `json:"relation,omitempty"`
	Rollup      *RollupMetadata   `json:"rollup,omitempty"`
}

// DatabaseQuery is used for quering a database.
type DatabaseQuery struct {
	Filter      *DatabaseQueryFilter `json:"filter,omitempty"`
	Sorts       []DatabaseQuerySort  `json:"sorts,omitempty"`
	StartCursor string               `json:"start_cursor,omitempty"`
	PageSize    int                  `json:"page_size,omitempty"`
}

// DatabaseQueryResponse contains the results and pagination data from a query request.
type DatabaseQueryResponse struct {
	Results    []Page  `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

// DatabaseQueryFilter is used to filter database contents.
// See: https://developers.notion.com/reference/post-database-query#post-database-query-filter
type DatabaseQueryFilter struct {
	Property string `json:"property,omitempty"`

	Text        *TextDatabaseQueryFilter        `json:"text,omitempty"`
	Number      *NumberDatabaseQueryFilter      `json:"number,omitempty"`
	Checkbox    *CheckboxDatabaseQueryFilter    `json:"checkbox,omitempty"`
	Select      *SelectDatabaseQueryFilter      `json:"select,omitempty"`
	MultiSelect *MultiSelectDatabaseQueryFilter `json:"multi_select,omitempty"`
	Date        *DateDatabaseQueryFilter        `json:"date,omitempty"`
	People      *PeopleDatabaseQueryFilter      `json:"people,omitempty"`
	Files       *FilesDatabaseQueryFilter       `json:"files,omitempty"`
	Relation    *RelationDatabaseQueryFilter    `json:"relation,omitempty"`

	Or  []DatabaseQueryFilter `json:"or,omitempty"`
	And []DatabaseQueryFilter `json:"and,omitempty"`
}

type TextDatabaseQueryFilter struct {
	Equals         string `json:"equals,omitempty"`
	DoesNotEqual   string `json:"does_not_equal,omitempty"`
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	StartsWith     string `json:"starts_with,omitempty"`
	EndsWith       string `json:"ends_with,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type NumberDatabaseQueryFilter struct {
	Equals               *int `json:"equals,omitempty"`
	DoesNotEqual         *int `json:"does_not_equal,omitempty"`
	GreaterThan          *int `json:"greater_than,omitempty"`
	LessThan             *int `json:"less_than,omitempty"`
	GreaterThanOrEqualTo *int `json:"greater_than_or_equal_to,omitempty"`
	LessThanOrEqualTo    *int `json:"less_than_or_equal_to,omitempty"`
	IsEmpty              bool `json:"is_empty,omitempty"`
	IsNotEmpty           bool `json:"is_not_empty,omitempty"`
}

type CheckboxDatabaseQueryFilter struct {
	Equals       *bool `json:"equals,omitempty"`
	DoesNotEqual *bool `json:"does_not_equal,omitempty"`
}

type SelectDatabaseQueryFilter struct {
	Equals       string `json:"equals,omitempty"`
	DoesNotEqual string `json:"does_not_equal,omitempty"`
	IsEmpty      bool   `json:"is_empty,omitempty"`
	IsNotEmpty   bool   `json:"is_not_empty,omitempty"`
}

type MultiSelectDatabaseQueryFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type DateDatabaseQueryFilter struct {
	Equals     *time.Time `json:"equals,omitempty"`
	Before     *time.Time `json:"before,omitempty"`
	After      *time.Time `json:"after,omitempty"`
	OnOrBefore *time.Time `json:"on_or_before,omitempty"`
	OnOrAfter  *time.Time `json:"on_or_after,omitempty"`
	IsEmpty    bool       `json:"is_empty,omitempty"`
	IsNotEmpty bool       `json:"is_not_empty,omitempty"`
	PastWeek   *struct{}  `json:"past_week,omitempty"`
	PastMonth  *struct{}  `json:"past_month,omitempty"`
	PastYear   *struct{}  `json:"past_year,omitempty"`
	NextWeek   *struct{}  `json:"next_week,omitempty"`
	NextMonth  *struct{}  `json:"next_month,omitempty"`
	NextYear   *struct{}  `json:"next_year,omitempty"`
}

type PeopleDatabaseQueryFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type FilesDatabaseQueryFilter struct {
	IsEmpty    bool `json:"is_empty,omitempty"`
	IsNotEmpty bool `json:"is_not_empty,omitempty"`
}

type RelationDatabaseQueryFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type FormulaDatabaseQueryFilter struct {
	Text     TextDatabaseQueryFilter     `json:"text,omitempty"`
	Checkbox CheckboxDatabaseQueryFilter `json:"checkbox,omitempty"`
	Number   NumberDatabaseQueryFilter   `json:"number,omitempty"`
	Date     DateDatabaseQueryFilter     `json:"date,omitempty"`
}

type DatabaseQuerySort struct {
	Property  string        `json:"property,omitempty"`
	Timestamp SortTimestamp `json:"timestamp,omitempty"`
	Direction SortDirection `json:"direction,omitempty"`
}

// CreateDatabaseParams are the params used for creating a database.
type CreateDatabaseParams struct {
	ParentPageID string
	Title        []RichText
	Properties   DatabaseProperties
}

type (
	DatabasePropertyType string
	NumberFormat         string
	FormulaResultType    string
	RollupResultType     string
	SortTimestamp        string
	SortDirection        string
)

const (
	// Database property type enums.
	DBPropTypeTitle          DatabasePropertyType = "title"
	DBPropTypeRichText       DatabasePropertyType = "rich_text"
	DBPropTypeNumber         DatabasePropertyType = "number"
	DBPropTypeSelect         DatabasePropertyType = "select"
	DBPropTypeMultiSelect    DatabasePropertyType = "multi_select"
	DBPropTypeDate           DatabasePropertyType = "date"
	DBPropTypePeople         DatabasePropertyType = "people"
	DBPropTypeFiles          DatabasePropertyType = "files"
	DBPropTypeCheckbox       DatabasePropertyType = "checkbox"
	DBPropTypeURL            DatabasePropertyType = "url"
	DBPropTypeEmail          DatabasePropertyType = "email"
	DBPropTypePhoneNumber    DatabasePropertyType = "phone_number"
	DBPropTypeFormula        DatabasePropertyType = "formula"
	DBPropTypeRelation       DatabasePropertyType = "relation"
	DBPropTypeRollup         DatabasePropertyType = "rollup"
	DBPropTypeCreatedTime    DatabasePropertyType = "created_time"
	DBPropTypeCreatedBy      DatabasePropertyType = "created_by"
	DBPropTypeLastEditedTime DatabasePropertyType = "last_edited_time"
	DBPropTypeLastEditedBy   DatabasePropertyType = "last_edited_by"

	// Number format enums.
	NumberFormatNumber           NumberFormat = "number"
	NumberFormatNumberWithCommas NumberFormat = "number_with_commas"
	NumberFormatPercent          NumberFormat = "percent"
	NumberFormatDollar           NumberFormat = "dollar"
	NumberFormatEuro             NumberFormat = "euro"
	NumberFormatPound            NumberFormat = "pound"
	NumberFormatPonud            NumberFormat = "yen"
	NumberFormatRuble            NumberFormat = "ruble"
	NumberFormatRupee            NumberFormat = "rupee"
	NumberFormatWon              NumberFormat = "won"
	NumberformatYuan             NumberFormat = "yuan"

	// Formula result type enums.
	FormulaResultTypeString  FormulaResultType = "string"
	FormulaResultTypeNumber  FormulaResultType = "number"
	FormulaResultTypeBoolean FormulaResultType = "boolean"
	FormulaResultTypeDate    FormulaResultType = "date"

	// Rollup result type enums.
	RollupResultTypeNumber RollupResultType = "number"
	RollupResultTypeDate   RollupResultType = "date"
	RollupResultTypeArray  RollupResultType = "array"

	// Sort timestamp enums.
	SortTimeStampCreatedTime    SortTimestamp = "created_time"
	SortTimeStampLastEditedTime SortTimestamp = "last_edited_time"

	// Sort direction enums.
	SortDirAsc  SortDirection = "ascending"
	SortDirDesc SortDirection = "descending"
)

// Metadata returns the underlying property metadata, based on its `type` field.
// When type is unknown/unmapped or doesn't have additional properies, `nil` is returned.
func (prop DatabaseProperty) Metadata() interface{} {
	switch prop.Type {
	case "title":
		return prop.Title
	case "number":
		return prop.Number
	case "select":
		return prop.Select
	case "multi_select":
		return prop.MultiSelect
	case "formula":
		return prop.Formula
	case "relation":
		return prop.Relation
	case "rollup":
		return prop.Rollup
	default:
		return nil
	}
}

// Value returns the underlying result value of an evaluated formula.
func (f FormulaResult) Value() interface{} {
	switch f.Type {
	case FormulaResultTypeString:
		return f.String
	case FormulaResultTypeNumber:
		return f.Number
	case FormulaResultTypeBoolean:
		return f.Boolean
	case FormulaResultTypeDate:
		return f.Date
	default:
		return nil
	}
}

// Value returns the underlying result value of an evaluated rollup.
func (r RollupResult) Value() interface{} {
	switch r.Type {
	case RollupResultTypeNumber:
		return r.Number
	case RollupResultTypeDate:
		return r.Date
	case RollupResultTypeArray:
		return r.Array
	default:
		return nil
	}
}

// Validate validates params for creating a database.
func (p CreateDatabaseParams) Validate() error {
	if p.ParentPageID == "" {
		return errors.New("parent page ID is required")
	}
	if p.Properties == nil {
		return errors.New("database properties are required")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (p CreateDatabaseParams) MarshalJSON() ([]byte, error) {
	type CreatePageParamsDTO struct {
		Parent     Parent             `json:"parent"`
		Title      []RichText         `json:"title,omitempty"`
		Properties DatabaseProperties `json:"properties"`
	}

	parent := Parent{
		Type:   ParentTypePage,
		PageID: p.ParentPageID,
	}

	dto := CreatePageParamsDTO{
		Parent:     parent,
		Title:      p.Title,
		Properties: p.Properties,
	}

	return json.Marshal(dto)
}
