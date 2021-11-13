package domain

type (
	FileStatus int
	FileType   string
)

const (
	ClientUploadInProgress FileStatus = iota
	UploadedByClient
	ClientUploadError
	StorageUploadInProgress
	UploadedToStorage
	StorageUploadError
)

const (
	Image    FileType = "image"
	Document FileType = "document"
	Other    FileType = "other"
)

type File struct {
	Object
	UserID      ID
	Type        FileType
	ContentType string
	Name        string
	Size        int64
	Status      FileStatus
	URL         string
}
