package notion

type Parent struct {
	Type ParentType `json:"type,omitempty"`

	PageID     string `json:"page_id,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
	Workspace  bool   `json:"workspace,omitempty"`
}

type ParentType string

const (
	ParentTypeDatabase  ParentType = "database_id"
	ParentTypePage      ParentType = "page_id"
	ParentTypeWorkspace ParentType = "workspace"
)
