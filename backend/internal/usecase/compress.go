package usecase

import (
	"fmt"
	"shucompress/internal/domain"
	"shucompress/pkg/utils"
	"shucompress/pkg/compressor"
)

type CompressUsecase struct {
	fileRepo    domain.FileRepository
	compressors []domain.CompressorService
}

func NewCompressUsecase(
	fileRepo domain.FileRepository,
	compressors ...domain.CompressorService,
) *CompressUsecase {
	return &CompressUsecase{
		fileRepo:    fileRepo,
		compressors: compressors,
	}
}

type CompressInput struct {
    FileData     []byte
    OriginalName string
    Quality      string
    TargetSize   int64  // 0 = no target
}

type CompressOutput struct {
	OutputPath     string
	Filename       string
	OriginalSize   int64
	CompressedSize int64
	Ratio          float64 // persentase pengurangan
}

func (uc *CompressUsecase) Execute(input CompressInput) (*CompressOutput, error) {
	// 1. Deteksi tipe file
	fileType, supported := utils.DetectFileType(input.OriginalName)
	if !supported {
		return nil, fmt.Errorf("format file tidak didukung")
	}

	// 2. Simpan file upload ke tmp
	tempPath, err := uc.fileRepo.Save(input.FileData, input.OriginalName)
	if err != nil {
		return nil, fmt.Errorf("gagal simpan file: %w", err)
	}
	defer uc.fileRepo.Delete(tempPath) // cleanup input setelah selesai

	// 3. Ambil size original
	originalSize, err := uc.fileRepo.GetSize(tempPath)
	if err != nil {
		return nil, fmt.Errorf("gagal baca ukuran file: %w", err)
	}

	// 4. Cari compressor yang support tipe ini
	var selectedCompressor domain.CompressorService
	for _, c := range uc.compressors {
		if c.Supports(fileType) {
			selectedCompressor = c
			break
		}
	}
	if selectedCompressor == nil {
		return nil, fmt.Errorf("tidak ada compressor untuk tipe: %s", fileType)
	}

	// 5. Jalankan compress
	req := domain.CompressRequest{
		OriginalName: input.OriginalName,
		TempPath:     tempPath,
		FileType:     fileType,
		Quality:      input.Quality,
		TargetSize:   input.TargetSize,
	}

	var result domain.CompressResult

	// Kalau ada target size dan file adalah PDF, pakai iterative compression
	if input.TargetSize > 0 && fileType == domain.FileTypePDF {
		pdfC, ok := selectedCompressor.(*compressor.PDFCompressor)
		if ok {
			result, err = pdfC.CompressToTarget(req, input.TargetSize)
		} else {
			result, err = selectedCompressor.Compress(req)
		}
	} else {
		result, err = selectedCompressor.Compress(req)
	}

	if err != nil {
		return nil, fmt.Errorf("gagal compress: %w", err)
	}

	// 6. Ambil size hasil
	compressedSize, err := uc.fileRepo.GetSize(result.OutputPath)
	if err != nil {
		uc.fileRepo.Delete(result.OutputPath)
		return nil, fmt.Errorf("gagal baca ukuran output: %w", err)
	}

	// 7. Kalau hasil lebih besar dari original, return error informatif
	if compressedSize >= originalSize {
		uc.fileRepo.Delete(result.OutputPath)
		return nil, fmt.Errorf("file sudah optimal, tidak bisa dikecilkan lebih jauh")
	}

	// 8. Hitung ratio pengurangan
	ratio := float64(originalSize-compressedSize) / float64(originalSize) * 100

	return &CompressOutput{
		OutputPath:     result.OutputPath,
		Filename:       result.Filename,
		OriginalSize:   originalSize,
		CompressedSize: compressedSize,
		Ratio:          ratio,
	}, nil
}