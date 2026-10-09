package models

type FileInfo struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	OriginalSize   int64  `json:"originalSize"`
	CompressedSize int64  `json:"compressedSize"`
	Type           string `json:"type"`
	Modified       string `json:"modified"`
	IsDir          bool   `json:"isDir"`
}

type ArchiveResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ProgressEvent struct {
	Percent     int    `json:"percent"`
	CurrentFile string `json:"currentFile"`
	Status      string `json:"status"`
}
