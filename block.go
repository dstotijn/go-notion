package notion

import (
	"encoding/json"
	"time"
)

// Block represents content on the Notion platform.
// See: https://developers.notion.com/reference/block
type Block struct {
	Object         string     `json:"object"`
	ID             string     `json:"id,omitempty"`
	Type           BlockType  `json:"type,omitempty"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Archived       *bool      `json:"archived,omitempty"`

	Paragraph        *RichTextBlock   `json:"paragraph,omitempty"`
	Heading1         *Heading         `json:"heading_1,omitempty"`
	Heading2         *Heading         `json:"heading_2,omitempty"`
	Heading3         *Heading         `json:"heading_3,omitempty"`
	BulletedListItem *RichTextBlock   `json:"bulleted_list_item,omitempty"`
	NumberedListItem *RichTextBlock   `json:"numbered_list_item,omitempty"`
	ToDo             *ToDo            `json:"to_do,omitempty"`
	Toggle           *RichTextBlock   `json:"toggle,omitempty"`
	ChildPage        *ChildPage       `json:"child_page,omitempty"`
	ChildDatabase    *ChildDatabase   `json:"child_database,omitempty"`
	Callout          *Callout         `json:"callout,omitempty"`
	Quote            *RichTextBlock   `json:"quote,omitempty"`
	Code             *Code            `json:"code,omitempty"`
	Embed            *Embed           `json:"embed,omitempty"`
	Image            *FileBlock       `json:"image,omitempty"`
	Video            *FileBlock       `json:"video,omitempty"`
	File             *FileBlock       `json:"file,omitempty"`
	PDF              *FileBlock       `json:"pdf,omitempty"`
	Bookmark         *Bookmark        `json:"bookmark,omitempty"`
	Equation         *Equation        `json:"equation,omitempty"`
	Divider          *Divider         `json:"divider,omitempty"`
	TableOfContents  *TableOfContents `json:"table_of_contents,omitempty"`
	Breadcrumb       *Breadcrumb      `json:"breadcrumb,omitempty"`
	ColumnList       *ColumnList      `json:"column_list,omitempty"`
	Column           *Column          `json:"column,omitempty"`
	LinkPreview      *LinkPreview     `json:"link_preview,omitempty"`
	LinkToPage       *LinkToPage      `json:"link_to_page,omitempty"`
	SyncedBlock      *SyncedBlock     `json:"synced_block,omitempty"`
	Template         *RichTextBlock   `json:"template,omitempty"`
}

type RichTextBlock struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

type Heading struct {
	Text []RichText `json:"text"`
}

type ToDo struct {
	RichTextBlock
	Checked *bool `json:"checked,omitempty"`
}

type ChildPage struct {
	Title string `json:"title"`
}

type ChildDatabase struct {
	Title string `json:"title"`
}

type Callout struct {
	RichTextBlock
	Icon *Icon `json:"icon,omitempty"`
}

type Code struct {
	RichTextBlock
	Language *string `json:"language,omitempty"`
}

type Embed struct {
	URL string `json:"url"`
}

type FileBlock struct {
	Type FileType `json:"type"`

	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
	Caption  []RichText    `json:"caption,omitempty"`
}

type Bookmark struct {
	URL     string     `json:"url"`
	Caption []RichText `json:"caption,omitempty"`
}

type ColumnList struct {
	Children []Block `json:"children,omitempty"`
}

type Column struct {
	Children []Block `json:"children,omitempty"`
}

type LinkToPage struct {
	Type LinkToPageType `json:"type"`

	PageID     string `json:"page_id,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
}

type LinkToPageType string

const (
	LinkToPageTypePageID     LinkToPageType = "page_id"
	LinkToPageTypeDatabaseID LinkToPageType = "database_id"
)

type SyncedBlock struct {
	SyncedFrom *SyncedFrom `json:"synced_from"`
	Children   []Block     `json:"children,omitempty"`
}

type SyncedFrom struct {
	Type    SyncedFromType `json:"type"`
	BlockID string         `json:"block_id"`
}

type SyncedFromType string

const SyncedFromTypeBlockID SyncedFromType = "block_id"

type (
	Divider         struct{}
	TableOfContents struct{}
	Breadcrumb      struct{}
)

type BlockType string

const (
	BlockTypeParagraph        BlockType = "paragraph"
	BlockTypeHeading1         BlockType = "heading_1"
	BlockTypeHeading2         BlockType = "heading_2"
	BlockTypeHeading3         BlockType = "heading_3"
	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"
	BlockTypeToDo             BlockType = "to_do"
	BlockTypeToggle           BlockType = "toggle"
	BlockTypeChildPage        BlockType = "child_page"
	BlockTypeChildDatabase    BlockType = "child_database"
	BlockTypeCallout          BlockType = "callout"
	BlockTypeQuote            BlockType = "quote"
	BlockTypeCode             BlockType = "code"
	BlockTypeEmbed            BlockType = "embed"
	BlockTypeImage            BlockType = "image"
	BlockTypeVideo            BlockType = "video"
	BlockTypeFile             BlockType = "file"
	BlockTypePDF              BlockType = "pdf"
	BlockTypeBookmark         BlockType = "bookmark"
	BlockTypeEquation         BlockType = "equation"
	BlockTypeDivider          BlockType = "divider"
	BlockTypeTableOfContents  BlockType = "table_of_contents"
	BlockTypeBreadCrumb       BlockType = "breadcrumb"
	BlockTypeColumnList       BlockType = "column_list"
	BlockTypeColumn           BlockType = "column"
	BlockTypeLinkPreview      BlockType = "link_preview"
	BlockTypeLinkToPage       BlockType = "link_to_page"
	BlockTypeSyncedBlock      BlockType = "synced_block"
	BlockTypeTemplate         BlockType = "template"
	BlockTypeUnsupported      BlockType = "unsupported"
)

type PaginationQuery struct {
	StartCursor string
	PageSize    int
}

// BlockChildrenResponse contains results (block children) and pagination data returned from a find request.
type BlockChildrenResponse struct {
	Results    []Block `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

// MarshalJSON implements json.Marshaler.
func (b Block) MarshalJSON() ([]byte, error) {
	type blockAlias Block

	alias := blockAlias(b)
	alias.Object = "block"

	return json.Marshal(alias)
}
