package http

import (
	"fmt"
	"os"
	"strconv"
	"shucompress/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type CompressHandler struct {
	compressUC *usecase.CompressUsecase
}

func NewCompressHandler(compressUC *usecase.CompressUsecase) *CompressHandler {
	return &CompressHandler{compressUC: compressUC}
}

func (h *CompressHandler) Compress(c *fiber.Ctx) error {
	// 1. Ambil file dari form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "File tidak ditemukan. Pastikan field name adalah 'file'",
		})
	}

	// 2. Baca isi file ke bytes
	f, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file"})
	}
	defer f.Close()

	fileData := make([]byte, file.Size)
	if _, err := f.Read(fileData); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca file"})
	}

	// 3. Ambil quality dari form (opsional)
	quality := c.FormValue("quality", "medium")

	// 4. Ambil target_size dalam MB, convert ke bytes
	targetSize := int64(0)
	if mb, err := strconv.ParseFloat(c.FormValue("target_size", "0"), 64); err == nil && mb > 0 {
		targetSize = int64(mb * 1024 * 1024)
	}

	// 5. Jalankan usecase
	output, err := h.compressUC.Execute(usecase.CompressInput{
		FileData:     fileData,
		OriginalName: file.Filename,
		Quality:      quality,
		TargetSize:   targetSize,
	})
	if err != nil {
		return c.Status(422).JSON(fiber.Map{"error": err.Error()})
	}

	// 6. Set header info size untuk frontend
	c.Set("X-Original-Size", fmt.Sprintf("%d", output.OriginalSize))
	c.Set("X-Compressed-Size", fmt.Sprintf("%d", output.CompressedSize))
	c.Set("X-Compression-Ratio", fmt.Sprintf("%.1f", output.Ratio))
	c.Set("Access-Control-Expose-Headers", "X-Original-Size, X-Compressed-Size, X-Compression-Ratio")

	// 7. Kirim file & cleanup
	defer os.Remove(output.OutputPath)
	return c.Download(output.OutputPath, output.Filename)
}