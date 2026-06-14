package compressor

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"shucompress/internal/domain"
	"strings"

	"golang.org/x/image/draw"
)

type ImageCompressor struct{}

func NewImageCompressor() *ImageCompressor {
	return &ImageCompressor{}
}

func (i *ImageCompressor) Supports(fileType domain.FileType) bool {
	return fileType == domain.FileTypeImage
}

func (i *ImageCompressor) Compress(req domain.CompressRequest) (domain.CompressResult, error) {
	// Buka file gambar
	inputFile, err := os.Open(req.TempPath)
	if err != nil {
		return domain.CompressResult{}, fmt.Errorf("gagal buka file: %w", err)
	}
	defer inputFile.Close()

	// Decode gambar (auto detect format)
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return domain.CompressResult{}, fmt.Errorf("gagal decode gambar: %w", err)
	}

	img = resizeIfNeeded(img, 1920)

	// Tentukan JPEG quality berdasarkan pilihan user
	quality := jpegQuality(req.Quality)

	// Output selalu JPEG (lebih kecil dari PNG)
	ext := strings.ToLower(filepath.Ext(req.OriginalName))
	outputExt := ".jpg"
	if ext == ".png" {
		outputExt = ".jpg" // PNG → JPEG untuk size lebih kecil
	}

	outputPath := req.TempPath + "_compressed" + outputExt

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return domain.CompressResult{}, fmt.Errorf("gagal buat output file: %w", err)
	}
	defer outputFile.Close()

	// Encode dengan quality yang dipilih
	if outputExt == ".jpg" {
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: quality})
	} else {
		err = png.Encode(outputFile, img)
	}

	if err != nil {
		return domain.CompressResult{}, fmt.Errorf("gagal encode gambar: %w", err)
	}

	// Ganti ekstensi di filename output
	originalBase := strings.TrimSuffix(req.OriginalName, filepath.Ext(req.OriginalName))

	return domain.CompressResult{
		OutputPath: outputPath,
		Filename:   "compressed_" + originalBase + outputExt,
	}, nil
}

// jpegQuality convert pilihan user ke angka JPEG quality
func jpegQuality(quality string) int {
	switch quality {
	case "low":
		return 40
	case "medium":
		return 65
	case "high":
		return 85
	default:
		return 65
	}
}


func resizeIfNeeded(img image.Image, maxWidth int) image.Image {
    bounds := img.Bounds()
    w := bounds.Dx()
    if w <= maxWidth {
        return img // sudah kecil, skip
    }
    h := bounds.Dy()
    newH := h * maxWidth / w
    dst := image.NewRGBA(image.Rect(0, 0, maxWidth, newH))
    draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
    return dst
}