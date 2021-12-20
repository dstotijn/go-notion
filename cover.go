package notion

import "errors"

type CoverType string

const (
	CoverTypeFile     CoverType = "file"
	CoverTypeExternal CoverType = "external"
)

type Cover struct {
	Type CoverType `json:"type"`

	File     *CoverFile     `json:"file,omitempty"`
	External *CoverExternal `json:"external,omitempty"`
}

type CoverFile struct {
	URL        string   `json:"url"`
	ExpiryTime DateTime `json:"expiry_time"`
}

type CoverExternal struct {
	URL string `json:"url"`
}

func (cover Cover) Validate() error {
	if cover.Type == "" {
		return errors.New("cover type cannot be empty")
	}

	if cover.Type == CoverTypeExternal && cover.External == nil {
		return errors.New("cover external cannot be empty")
	}

	return nil
}
