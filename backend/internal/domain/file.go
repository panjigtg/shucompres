package domain

// FileType represents supported file types
type FileType string

const (
	FileTypePDF   FileType = "pdf"
	FileTypeImage FileType = "image"
	FileTypePPTX  FileType = "pptx"
)

// CompressRequest adalah data yang masuk dari user
type CompressRequest struct {
    OriginalName string
    TempPath     string
    FileType     FileType
    Quality      string
    TargetSize   int64  // 0 = tidak ada target, >0 = target bytes
}

// CompressResult adalah hasil kompresi
type CompressResult struct {
	OutputPath     string
	OriginalSize   int64
	CompressedSize int64
	Filename       string
}

// FileRepository interface — cara simpan & baca file sementara
type FileRepository interface {
	Save(data []byte, filename string) (path string, err error)
	Delete(path string) error
	GetSize(path string) (int64, error)
}

// CompressorService interface — kontrak untuk semua compressor
type CompressorService interface {
	Compress(req CompressRequest) (CompressResult, error)
	Supports(fileType FileType) bool
}