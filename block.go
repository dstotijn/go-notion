package notion

import (
	"encoding/json"
	"fmt"
	"time"
)

// Block represents content on the Notion platform.
// See: https://developers.notion.com/reference/block
type Block interface {
	ID() string
	Parent() Parent
	CreatedTime() time.Time
	CreatedBy() BaseUser
	LastEditedBy() BaseUser
	LastEditedTime() time.Time
	HasChildren() bool
	Archived() bool
	json.Marshaler
}

type blockDTO struct {
	ID             string     `json:"id,omitempty"`
	Parent         *Parent    `json:"parent,omitempty"`
	Type           BlockType  `json:"type,omitempty"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	CreatedBy      *BaseUser  `json:"created_by,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	LastEditedBy   *BaseUser  `json:"last_edited_by,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Archived       *bool      `json:"archived,omitempty"`

	Paragraph        *ParagraphBlock        `json:"paragraph,omitempty"`
	Heading1         *Heading1Block         `json:"heading_1,omitempty"`
	Heading2         *Heading2Block         `json:"heading_2,omitempty"`
	Heading3         *Heading3Block         `json:"heading_3,omitempty"`
	BulletedListItem *BulletedListItemBlock `json:"bulleted_list_item,omitempty"`
	NumberedListItem *NumberedListItemBlock `json:"numbered_list_item,omitempty"`
	ToDo             *ToDoBlock             `json:"to_do,omitempty"`
	Toggle           *ToggleBlock           `json:"toggle,omitempty"`
	ChildPage        *ChildPageBlock        `json:"child_page,omitempty"`
	ChildDatabase    *ChildDatabaseBlock    `json:"child_database,omitempty"`
	Callout          *CalloutBlock          `json:"callout,omitempty"`
	Quote            *QuoteBlock            `json:"quote,omitempty"`
	Code             *CodeBlock             `json:"code,omitempty"`
	Embed            *EmbedBlock            `json:"embed,omitempty"`
	Image            *ImageBlock            `json:"image,omitempty"`
	Audio            *AudioBlock            `json:"audio,omitempty"`
	Video            *VideoBlock            `json:"video,omitempty"`
	File             *FileBlock             `json:"file,omitempty"`
	PDF              *PDFBlock              `json:"pdf,omitempty"`
	Bookmark         *BookmarkBlock         `json:"bookmark,omitempty"`
	Equation         *EquationBlock         `json:"equation,omitempty"`
	Divider          *DividerBlock          `json:"divider,omitempty"`
	TableOfContents  *TableOfContentsBlock  `json:"table_of_contents,omitempty"`
	Breadcrumb       *BreadcrumbBlock       `json:"breadcrumb,omitempty"`
	ColumnList       *ColumnListBlock       `json:"column_list,omitempty"`
	Column           *ColumnBlock           `json:"column,omitempty"`
	Table            *TableBlock            `json:"table,omitempty"`
	TableRow         *TableRowBlock         `json:"table_row,omitempty"`
	LinkPreview      *LinkPreviewBlock      `json:"link_preview,omitempty"`
	LinkToPage       *LinkToPageBlock       `json:"link_to_page,omitempty"`
	SyncedBlock      *SyncedBlock           `json:"synced_block,omitempty"`
	Template         *TemplateBlock         `json:"template,omitempty"`
}

type baseBlock struct {
	id             string
	parent         Parent
	createdTime    time.Time
	createdBy      BaseUser
	lastEditedTime time.Time
	lastEditedBy   BaseUser
	hasChildren    bool
	archived       bool
}

// ID returns the identifier (UUIDv4) for the block.
func (b baseBlock) ID() string {
	return b.id
}

func (b baseBlock) CreatedTime() time.Time {
	return b.createdTime
}

func (b baseBlock) CreatedBy() BaseUser {
	return b.createdBy
}

func (b baseBlock) LastEditedTime() time.Time {
	return b.lastEditedTime
}

func (b baseBlock) LastEditedBy() BaseUser {
	return b.lastEditedBy
}

func (b baseBlock) HasChildren() bool {
	return b.hasChildren
}

func (b baseBlock) Archived() bool {
	return b.archived
}

func (b baseBlock) Parent() Parent {
	return b.parent
}

type ParagraphBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b ParagraphBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ParagraphBlock
		dto        struct {
			Paragraph blockAlias `json:"paragraph"`
		}
	)

	return json.Marshal(dto{
		Paragraph: blockAlias(b),
	})
}

type BulletedListItemBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b BulletedListItemBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias BulletedListItemBlock
		dto        struct {
			BulletedListItem blockAlias `json:"bulleted_list_item"`
		}
	)

	return json.Marshal(dto{
		BulletedListItem: blockAlias(b),
	})
}

type NumberedListItemBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b NumberedListItemBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias NumberedListItemBlock
		dto        struct {
			NumberedListItem blockAlias `json:"numbered_list_item"`
		}
	)

	return json.Marshal(dto{
		NumberedListItem: blockAlias(b),
	})
}

type QuoteBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b QuoteBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias QuoteBlock
		dto        struct {
			Quote blockAlias `json:"quote"`
		}
	)

	return json.Marshal(dto{
		Quote: blockAlias(b),
	})
}

type ToggleBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b ToggleBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ToggleBlock
		dto        struct {
			Toggle blockAlias `json:"toggle"`
		}
	)

	return json.Marshal(dto{
		Toggle: blockAlias(b),
	})
}

type TemplateBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b TemplateBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias TemplateBlock
		dto        struct {
			Template blockAlias `json:"template"`
		}
	)

	return json.Marshal(dto{
		Template: blockAlias(b),
	})
}

type Heading1Block struct {
	baseBlock

	RichText     []RichText `json:"rich_text"`
	Children     []Block    `json:"children,omitempty"`
	Color        Color      `json:"color,omitempty"`
	IsToggleable bool       `json:"is_toggleable"`
}

// MarshalJSON implements json.Marshaler.
func (b Heading1Block) MarshalJSON() ([]byte, error) {
	type (
		blockAlias Heading1Block
		dto        struct {
			Heading1 blockAlias `json:"heading_1"`
		}
	)

	return json.Marshal(dto{
		Heading1: blockAlias(b),
	})
}

type Heading2Block struct {
	baseBlock

	RichText     []RichText `json:"rich_text"`
	Children     []Block    `json:"children,omitempty"`
	Color        Color      `json:"color,omitempty"`
	IsToggleable bool       `json:"is_toggleable"`
}

// MarshalJSON implements json.Marshaler.
func (b Heading2Block) MarshalJSON() ([]byte, error) {
	type (
		blockAlias Heading2Block
		dto        struct {
			Heading2 blockAlias `json:"heading_2"`
		}
	)

	return json.Marshal(dto{
		Heading2: blockAlias(b),
	})
}

type Heading3Block struct {
	baseBlock

	RichText     []RichText `json:"rich_text"`
	Children     []Block    `json:"children,omitempty"`
	Color        Color      `json:"color,omitempty"`
	IsToggleable bool       `json:"is_toggleable"`
}

// MarshalJSON implements json.Marshaler.
func (b Heading3Block) MarshalJSON() ([]byte, error) {
	type (
		blockAlias Heading3Block
		dto        struct {
			Heading3 blockAlias `json:"heading_3"`
		}
	)

	return json.Marshal(dto{
		Heading3: blockAlias(b),
	})
}

type ToDoBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Checked  *bool      `json:"checked,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b ToDoBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ToDoBlock
		dto        struct {
			ToDo blockAlias `json:"to_do"`
		}
	)

	return json.Marshal(dto{
		ToDo: blockAlias(b),
	})
}

type ChildPageBlock struct {
	baseBlock

	Title string `json:"title"`
}

// MarshalJSON implements json.Marshaler.
func (b ChildPageBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ChildPageBlock
		dto        struct {
			ChildPage blockAlias `json:"child_page"`
		}
	)

	return json.Marshal(dto{
		ChildPage: blockAlias(b),
	})
}

type ChildDatabaseBlock struct {
	baseBlock

	Title string `json:"title"`
}

// MarshalJSON implements json.Marshaler.
func (b ChildDatabaseBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ChildDatabaseBlock
		dto        struct {
			ChildDatabase blockAlias `json:"child_database"`
		}
	)

	return json.Marshal(dto{
		ChildDatabase: blockAlias(b),
	})
}

type CalloutBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Icon     *Icon      `json:"icon,omitempty"`
	Color    Color      `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b CalloutBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias CalloutBlock
		dto        struct {
			Callout blockAlias `json:"callout"`
		}
	)

	return json.Marshal(dto{
		Callout: blockAlias(b),
	})
}

type CodeBlock struct {
	baseBlock

	RichText []RichText `json:"rich_text"`
	Children []Block    `json:"children,omitempty"`
	Caption  []RichText `json:"caption,omitempty"`
	Language *string    `json:"language,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b CodeBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias CodeBlock
		dto        struct {
			Code blockAlias `json:"code"`
		}
	)

	return json.Marshal(dto{
		Code: blockAlias(b),
	})
}

type EmbedBlock struct {
	baseBlock

	URL string `json:"url"`
}

// MarshalJSON implements json.Marshaler.
func (b EmbedBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias EmbedBlock
		dto        struct {
			Embed blockAlias `json:"embed"`
		}
	)

	return json.Marshal(dto{
		Embed: blockAlias(b),
	})
}

type ImageBlock struct {
	baseBlock

	Type     FileType      `json:"type"`
	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
	Caption  []RichText    `json:"caption,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b ImageBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ImageBlock
		dto        struct {
			Image blockAlias `json:"image"`
		}
	)

	return json.Marshal(dto{
		Image: blockAlias(b),
	})
}

type AudioBlock struct {
	baseBlock

	Type     FileType      `json:"type"`
	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
	Caption  []RichText    `json:"caption,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b AudioBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ImageBlock
		dto        struct {
			Audio blockAlias `json:"audio"`
		}
	)

	return json.Marshal(dto{
		Audio: blockAlias(b),
	})
}

type VideoBlock struct {
	baseBlock

	Type     FileType      `json:"type"`
	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
	Caption  []RichText    `json:"caption,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b VideoBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias VideoBlock
		dto        struct {
			Video blockAlias `json:"video"`
		}
	)

	return json.Marshal(dto{
		Video: blockAlias(b),
	})
}

type FileBlock struct {
	baseBlock

	Type     FileType      `json:"type"`
	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
	Caption  []RichText    `json:"caption,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b FileBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias FileBlock
		dto        struct {
			File blockAlias `json:"file"`
		}
	)

	return json.Marshal(dto{
		File: blockAlias(b),
	})
}

type PDFBlock struct {
	baseBlock

	Type     FileType      `json:"type"`
	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
	Caption  []RichText    `json:"caption,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b PDFBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias PDFBlock
		dto        struct {
			PDF blockAlias `json:"pdf"`
		}
	)

	return json.Marshal(dto{
		PDF: blockAlias(b),
	})
}

type BookmarkBlock struct {
	baseBlock

	URL     string     `json:"url"`
	Caption []RichText `json:"caption,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b BookmarkBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias BookmarkBlock
		dto        struct {
			Bookmark blockAlias `json:"bookmark"`
		}
	)

	return json.Marshal(dto{
		Bookmark: blockAlias(b),
	})
}

type EquationBlock struct {
	baseBlock

	Expression string `json:"expression"`
}

// MarshalJSON implements json.Marshaler.
func (b EquationBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias EquationBlock
		dto        struct {
			Equation blockAlias `json:"equation"`
		}
	)

	return json.Marshal(dto{
		Equation: blockAlias(b),
	})
}

type ColumnListBlock struct {
	baseBlock

	Children []ColumnBlock `json:"children,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b ColumnListBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ColumnListBlock
		dto        struct {
			ColumnList blockAlias `json:"column_list"`
		}
	)

	return json.Marshal(dto{
		ColumnList: blockAlias(b),
	})
}

type ColumnBlock struct {
	baseBlock

	Children []Block `json:"children,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b ColumnBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias ColumnBlock
		dto        struct {
			Column blockAlias `json:"column"`
		}
	)

	return json.Marshal(dto{
		Column: blockAlias(b),
	})
}

type TableBlock struct {
	baseBlock

	TableWidth      int     `json:"table_width"`
	HasColumnHeader bool    `json:"has_column_header"`
	HasRowHeader    bool    `json:"has_row_header"`
	Children        []Block `json:"children,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b TableBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias TableBlock
		dto        struct {
			Table blockAlias `json:"table"`
		}
	)

	return json.Marshal(dto{
		Table: blockAlias(b),
	})
}

type TableRowBlock struct {
	baseBlock

	Cells [][]RichText `json:"cells"`
}

// MarshalJSON implements json.Marshaler.
func (b TableRowBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias TableRowBlock
		dto        struct {
			TableRow blockAlias `json:"table_row"`
		}
	)

	return json.Marshal(dto{
		TableRow: blockAlias(b),
	})
}

type LinkPreviewBlock struct {
	baseBlock

	URL string `json:"url"`
}

// MarshalJSON implements json.Marshaler.
func (b LinkPreviewBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias LinkPreviewBlock
		dto        struct {
			LinkPreview blockAlias `json:"link_preview"`
		}
	)

	return json.Marshal(dto{
		LinkPreview: blockAlias(b),
	})
}

type LinkToPageBlock struct {
	baseBlock

	Type       LinkToPageType `json:"type"`
	PageID     string         `json:"page_id,omitempty"`
	DatabaseID string         `json:"database_id,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b LinkToPageBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias LinkToPageBlock
		dto        struct {
			LinkToPage blockAlias `json:"link_to_page"`
		}
	)

	return json.Marshal(dto{
		LinkToPage: blockAlias(b),
	})
}

type LinkToPageType string

const (
	LinkToPageTypePageID     LinkToPageType = "page_id"
	LinkToPageTypeDatabaseID LinkToPageType = "database_id"
)

type SyncedBlock struct {
	baseBlock

	SyncedFrom *SyncedFrom `json:"synced_from"`
	Children   []Block     `json:"children,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b SyncedBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias SyncedBlock
		dto        struct {
			SyncedBlock blockAlias `json:"synced_block"`
		}
	)

	return json.Marshal(dto{
		SyncedBlock: blockAlias(b),
	})
}

type SyncedFrom struct {
	Type    SyncedFromType `json:"type"`
	BlockID string         `json:"block_id"`
}

type SyncedFromType string

const SyncedFromTypeBlockID SyncedFromType = "block_id"

type DividerBlock struct {
	baseBlock
}

// MarshalJSON implements json.Marshaler.
func (b DividerBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias DividerBlock
		dto        struct {
			Divider blockAlias `json:"divider"`
		}
	)

	return json.Marshal(dto{
		Divider: blockAlias(b),
	})
}

type TableOfContentsBlock struct {
	baseBlock

	Color Color `json:"color,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b TableOfContentsBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias TableOfContentsBlock
		dto        struct {
			TableOfContents blockAlias `json:"table_of_contents"`
		}
	)

	return json.Marshal(dto{
		TableOfContents: blockAlias(b),
	})
}

type BreadcrumbBlock struct {
	baseBlock
}

// MarshalJSON implements json.Marshaler.
func (b BreadcrumbBlock) MarshalJSON() ([]byte, error) {
	type (
		blockAlias BreadcrumbBlock
		dto        struct {
			Breadcrumb blockAlias `json:"breadcrumb"`
		}
	)

	return json.Marshal(dto{
		Breadcrumb: blockAlias(b),
	})
}

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
	BlockTypeAudio            BlockType = "audio"
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
	BlockTypeTable            BlockType = "table"
	BlockTypeTableRow         BlockType = "table_row"
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
	Results    []Block
	HasMore    bool
	NextCursor *string
}

func (resp *BlockChildrenResponse) UnmarshalJSON(b []byte) error {
	type responseDTO struct {
		Results    []blockDTO `json:"results"`
		HasMore    bool       `json:"has_more"`
		NextCursor *string    `json:"next_cursor"`
	}

	var dto responseDTO

	if err := json.Unmarshal(b, &dto); err != nil {
		return err
	}

	resp.HasMore = dto.HasMore
	resp.NextCursor = dto.NextCursor
	resp.Results = make([]Block, len(dto.Results))

	for i, blockDTO := range dto.Results {
		resp.Results[i] = blockDTO.Block()
	}

	return nil
}

func (dto blockDTO) Block() Block {
	baseBlock := baseBlock{
		id:          dto.ID,
		hasChildren: dto.HasChildren,
	}

	if dto.Parent != nil {
		baseBlock.parent = *dto.Parent
	}

	if dto.CreatedTime != nil {
		baseBlock.createdTime = *dto.CreatedTime
	}

	if dto.CreatedBy != nil {
		baseBlock.createdBy = *dto.CreatedBy
	}

	if dto.LastEditedTime != nil {
		baseBlock.lastEditedTime = *dto.LastEditedTime
	}

	if dto.LastEditedBy != nil {
		baseBlock.lastEditedBy = *dto.LastEditedBy
	}

	if dto.Archived != nil {
		baseBlock.archived = *dto.Archived
	}

	switch dto.Type {
	case BlockTypeParagraph:
		dto.Paragraph.baseBlock = baseBlock
		return dto.Paragraph
	case BlockTypeHeading1:
		dto.Heading1.baseBlock = baseBlock
		return dto.Heading1
	case BlockTypeHeading2:
		dto.Heading2.baseBlock = baseBlock
		return dto.Heading2
	case BlockTypeHeading3:
		dto.Heading3.baseBlock = baseBlock
		return dto.Heading3
	case BlockTypeBulletedListItem:
		dto.BulletedListItem.baseBlock = baseBlock
		return dto.BulletedListItem
	case BlockTypeNumberedListItem:
		dto.NumberedListItem.baseBlock = baseBlock
		return dto.NumberedListItem
	case BlockTypeToDo:
		dto.ToDo.baseBlock = baseBlock
		return dto.ToDo
	case BlockTypeToggle:
		dto.Toggle.baseBlock = baseBlock
		return dto.Toggle
	case BlockTypeChildPage:
		dto.ChildPage.baseBlock = baseBlock
		return dto.ChildPage
	case BlockTypeChildDatabase:
		dto.ChildDatabase.baseBlock = baseBlock
		return dto.ChildDatabase
	case BlockTypeCallout:
		dto.Callout.baseBlock = baseBlock
		return dto.Callout
	case BlockTypeQuote:
		dto.Quote.baseBlock = baseBlock
		return dto.Quote
	case BlockTypeCode:
		dto.Code.baseBlock = baseBlock
		return dto.Code
	case BlockTypeEmbed:
		dto.Embed.baseBlock = baseBlock
		return dto.Embed
	case BlockTypeImage:
		dto.Image.baseBlock = baseBlock
		return dto.Image
	case BlockTypeAudio:
		dto.Audio.baseBlock = baseBlock
		return dto.Audio
	case BlockTypeVideo:
		dto.Video.baseBlock = baseBlock
		return dto.Video
	case BlockTypeFile:
		dto.File.baseBlock = baseBlock
		return dto.File
	case BlockTypePDF:
		dto.PDF.baseBlock = baseBlock
		return dto.PDF
	case BlockTypeBookmark:
		dto.Bookmark.baseBlock = baseBlock
		return dto.Bookmark
	case BlockTypeEquation:
		dto.Equation.baseBlock = baseBlock
		return dto.Equation
	case BlockTypeDivider:
		dto.Divider.baseBlock = baseBlock
		return dto.Divider
	case BlockTypeTableOfContents:
		dto.TableOfContents.baseBlock = baseBlock
		return dto.TableOfContents
	case BlockTypeBreadCrumb:
		dto.Breadcrumb.baseBlock = baseBlock
		return dto.Breadcrumb
	case BlockTypeColumnList:
		dto.ColumnList.baseBlock = baseBlock
		return dto.ColumnList
	case BlockTypeColumn:
		dto.Column.baseBlock = baseBlock
		return dto.Column
	case BlockTypeTable:
		dto.Table.baseBlock = baseBlock
		return dto.Table
	case BlockTypeTableRow:
		dto.TableRow.baseBlock = baseBlock
		return dto.TableRow
	case BlockTypeLinkPreview:
		dto.LinkPreview.baseBlock = baseBlock
		return dto.LinkPreview
	case BlockTypeLinkToPage:
		dto.LinkToPage.baseBlock = baseBlock
		return dto.LinkToPage
	case BlockTypeSyncedBlock:
		dto.SyncedBlock.baseBlock = baseBlock
		return dto.SyncedBlock
	case BlockTypeTemplate:
		dto.Template.baseBlock = baseBlock
		return dto.Template
	default:
		panic(fmt.Sprintf("type %q is unsupported", dto.Type))
	}
}
