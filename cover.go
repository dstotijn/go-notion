package notion

import "errors"

type Cover struct {
	Type FileType `json:"type"`

	File     *FileFile     `json:"file,omitempty"`
	External *FileExternal `json:"external,omitempty"`
}

func (cover Cover) Validate() error {
	if cover.Type == "" {
		return errors.New("cover type cannot be empty")
	}

	if cover.Type == FileTypeExternal && cover.External == nil {
		return errors.New("cover external cannot be empty")
	}

	return nil
}
