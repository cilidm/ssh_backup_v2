package model

type FileInfo struct {
	FileName   string     `json:"file_name"`
	FileSource string     `json:"file_source"`
	FileTarget string     `json:"file_target"`
	FileSize   int64      `json:"file_size"`
	MD5        string     `json:"md5"`
	Status     FileStatus `json:"status"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"update_at"`
}

type FileInfoJson struct {
	FileName   string     `json:"file_name"`
	FileSource string     `json:"file_source"`
	FileTarget string     `json:"file_target"`
	FileSize   string     `json:"file_size"`
	Status     FileStatus `json:"status"`
	CreatedAt  string     `json:"created_at"`
}

type FileStatus int

const (
	Prossing  FileStatus = 1
	Processed FileStatus = 2
	ProssErr  FileStatus = 3
)
