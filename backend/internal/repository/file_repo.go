package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileRepository struct {
	BasePath string
}

func NewFileRepository(basePath string) *FileRepository {
	return &FileRepository{BasePath: basePath}
}

// Save menyimpan raw bytes ke tmp folder, return path-nya
func (r *FileRepository) Save(data []byte, filename string) (string, error) {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%d%s", timestamp, ext)
	fullPath := filepath.Join(r.BasePath, uniqueName)

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %w", err)
	}

	return fullPath, nil
}

// Delete hapus file dari tmp
func (r *FileRepository) Delete(path string) error {
	return os.Remove(path)
}

// GetSize return ukuran file dalam bytes
func (r *FileRepository) GetSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("gagal baca info file: %w", err)
	}
	return info.Size(), nil
}