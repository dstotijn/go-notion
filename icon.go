package notion

import "errors"

type IconType string

const (
	IconTypeEmoji    IconType = "emoji"
	IconTypeFile     IconType = "file"
	IconTypeExternal IconType = "external"
)

// Icon has one non-nil Emoji or External field, denoted by the corresponding
// IconType.
type Icon struct {
	Type IconType `json:"type"`

	Emoji    *string       `json:"emoji,omitempty"`
	File     *IconFile     `json:"file,omitempty"`
	External *IconExternal `json:"external,omitempty"`
}

type IconFile struct {
	URL        string   `json:"url"`
	ExpiryTime DateTime `json:"expiry_time"`
}

type IconExternal struct {
	URL string `json:"url"`
}

func (icon Icon) Validate() error {
	if icon.Type == "" {
		return errors.New("icon type cannot be empty")
	}

	if icon.Type == IconTypeEmoji && icon.Emoji == nil {
		return errors.New("icon emoji cannot be empty")
	}
	if icon.Type == IconTypeExternal && icon.External == nil {
		return errors.New("icon external cannot be empty")
	}

	return nil
}
