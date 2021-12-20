package notion

type FileFile struct {
	URL        string   `json:"url"`
	ExpiryTime DateTime `json:"expiry_time"`
}

type FileExternal struct {
	URL string `json:"url"`
}

type FileType string

const (
	FileTypeFile     FileType = "file"
	FileTypeExternal FileType = "external"
)
