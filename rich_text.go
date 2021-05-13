package notion

import "time"

type RichText struct {
	PlainText   string      `json:"plain_text"`
	HRef        *string     `json:"href"`
	Annotations Annotations `json:"annotations"`
	Type        string      `json:"type"`

	Text     *Text     `json:"text"`
	Mention  *Mention  `json:"mention"`
	Equation *Equation `json:"equation"`
}

type Equation struct {
	Expression string `json:"expression"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type Mention struct {
	Type string `json:"type"`

	User     *User            `json:"user"`
	Page     *PageMention     `json:"page"`
	Database *DatabaseMention `json:"database"`
	Date     *Date            `json:"date"`
}

type Date struct {
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end"`
}

type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link"`
}

type Link struct {
	URL string `json:"url"`
}

type PageMention struct {
	ID string `json:"id"`
}

type DatabaseMention struct {
	ID string `json:"id"`
}
