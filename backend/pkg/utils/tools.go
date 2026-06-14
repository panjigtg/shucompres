package utils

import (
	"fmt"
	"path/filepath"
	"shucompress/internal/domain"
	"strings"
)

// DetectFileType deteksi tipe file dari ekstensi
func DetectFileType(filename string) (domain.FileType, bool) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".pdf":
		return domain.FileTypePDF, true
	case ".jpg", ".jpeg", ".png", ".webp":
		return domain.FileTypeImage, true
	case ".ppt", ".pptx":
		return domain.FileTypePPTX, true
	default:
		return "", false
	}
}

// FormatBytes convert bytes ke human readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strings.TrimRight(strings.TrimRight(
			strings.Replace(fmt.Sprintf("%.1f", float64(bytes)), ",", ".", 1),
			"0"), ".") + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB"}
	return strings.TrimRight(strings.TrimRight(
		strings.Replace(fmt.Sprintf("%.1f", float64(bytes)/float64(div)), ",", ".", 1),
		"0"), ".") + " " + units[exp]
}