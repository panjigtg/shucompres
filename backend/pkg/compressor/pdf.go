package compressor

import (
	"fmt"
	"os/exec"
	"runtime"
	"shucompress/internal/domain"
	"os"
)

type PDFCompressor struct{}

func NewPDFCompressor() *PDFCompressor {
	return &PDFCompressor{}
}

func (p *PDFCompressor) Supports(fileType domain.FileType) bool {
	return fileType == domain.FileTypePDF
}

func (p *PDFCompressor) Compress(req domain.CompressRequest) (domain.CompressResult, error) {
	if !GhostscriptAvailable() {
		return domain.CompressResult{}, fmt.Errorf(
			"PDF compression membutuhkan Ghostscript. Install Ghostscript atau gunakan kompresi image/PPTX saja",
		)
	}

	outputPath := req.TempPath + "_compressed.pdf"

	quality := req.Quality
	validQualities := map[string]bool{
		"screen":  true,
		"ebook":   true,
		"printer": true,
	}
	if !validQualities[quality] {
		quality = "ebook"
	}

	gsCmd := ghostscriptBinary()

	cmd := exec.Command(
		gsCmd,
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.7",       // format lebih modern
		"-dNOPAUSE", "-dQUIET", "-dBATCH",

		// Image compression
		"-dAutoFilterColorImages=false",
		"-dColorImageFilter=/DCTEncode",   // JPEG untuk foto
		"-dColorImageResolution=150",      // 150dpi cukup untuk screen
		"-dGrayImageFilter=/DCTEncode",
		"-dGrayImageResolution=150",

		// Monochrome (teks hitam putih)
		"-dMonoImageFilter=/CCITTFaxEncode", // lossless untuk teks
		"-dMonoImageResolution=300",          // teks tetap tajam 300dpi

		// Jangan downsample kalau resolusi sudah rendah
		"-dDownsampleColorImages=true",
		"-dDownsampleGrayImages=true",
		"-dDownsampleMonoImages=false",    // jangan sentuh teks

		// Font
		"-dEmbedAllFonts=true",
		"-dSubsetFonts=true",              // embed hanya karakter yang dipakai

		// Hapus bloat
		"-dDetectDuplicateImages=true",    // dedupe gambar identik
		"-dCompressFonts=true",
		"-dOptimize=true",

		fmt.Sprintf("-sOutputFile=%s", outputPath),
		req.TempPath,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return domain.CompressResult{}, fmt.Errorf("ghostscript error: %w\n%s", err, output)
	}

	return domain.CompressResult{
		OutputPath: outputPath,
		Filename:   "compressed_" + req.OriginalName,
	}, nil
}

// ghostscriptBinary return nama binary sesuai OS
func ghostscriptBinary() string {
	if runtime.GOOS == "windows" {
		return "gswin64c"
	}
	return "gs"
}

func (p *PDFCompressor) CompressToTarget(req domain.CompressRequest, targetBytes int64) (domain.CompressResult, error) {
	if !GhostscriptAvailable() {
		return domain.CompressResult{}, fmt.Errorf(
			"PDF compression membutuhkan Ghostscript. Install Ghostscript atau gunakan target size yang lebih longgar",
		)
	}

    // Urutan quality yang dicoba, dari paling ringan ke paling agresif
    qualities := []string{"printer", "ebook", "screen"}

    // Tambah custom quality dengan DPI makin kecil
    customDPI := []int{120, 96, 72}

    gsCmd := ghostscriptBinary()

    for _, q := range qualities {
        outputPath := req.TempPath + "_compressed.pdf"

        cmd := exec.Command(
            gsCmd,
            "-sDEVICE=pdfwrite",
            "-dCompatibilityLevel=1.7",
            "-dNOPAUSE", "-dQUIET", "-dBATCH",
            "-dAutoFilterColorImages=false",
            "-dColorImageFilter=/DCTEncode",
            "-dGrayImageFilter=/DCTEncode",
            "-dMonoImageFilter=/CCITTFaxEncode",
            "-dMonoImageResolution=300",
            "-dDownsampleColorImages=true",
            "-dDownsampleGrayImages=true",
            "-dDetectDuplicateImages=true",
            "-dEmbedAllFonts=true",
            "-dSubsetFonts=true",
            "-dCompressFonts=true",
            "-dOptimize=true",
            fmt.Sprintf("-dPDFSETTINGS=/%s", q),
            fmt.Sprintf("-sOutputFile=%s", outputPath),
            req.TempPath,
        )

        if err := cmd.Run(); err != nil {
            os.Remove(outputPath)
            continue
        }

        info, err := os.Stat(outputPath)
        if err != nil {
            os.Remove(outputPath)
            continue
        }

        // Kalau sudah di bawah target, selesai
        if info.Size() <= targetBytes {
            return domain.CompressResult{
                OutputPath: outputPath,
                Filename:   "compressed_" + req.OriginalName,
            }, nil
        }

        os.Remove(outputPath)
    }

    // Kalau masih belum cukup kecil, coba DPI makin kecil
    for _, dpi := range customDPI {
        outputPath := req.TempPath + "_compressed.pdf"

        cmd := exec.Command(
            gsCmd,
            "-sDEVICE=pdfwrite",
            "-dCompatibilityLevel=1.7",
            "-dNOPAUSE", "-dQUIET", "-dBATCH",
            "-dAutoFilterColorImages=false",
            "-dColorImageFilter=/DCTEncode",
            fmt.Sprintf("-dColorImageResolution=%d", dpi),
            "-dGrayImageFilter=/DCTEncode",
            fmt.Sprintf("-dGrayImageResolution=%d", dpi),
            "-dMonoImageFilter=/CCITTFaxEncode",
            "-dMonoImageResolution=300",
            "-dDownsampleColorImages=true",
            "-dDownsampleGrayImages=true",
            "-dDetectDuplicateImages=true",
            "-dSubsetFonts=true",
            "-dCompressFonts=true",
            "-dOptimize=true",
            fmt.Sprintf("-sOutputFile=%s", outputPath),
            req.TempPath,
        )

        if err := cmd.Run(); err != nil {
            os.Remove(outputPath)
            continue
        }

        info, err := os.Stat(outputPath)
        if err != nil {
            os.Remove(outputPath)
            continue
        }

        if info.Size() <= targetBytes {
            return domain.CompressResult{
                OutputPath: outputPath,
                Filename:   "compressed_" + req.OriginalName,
            }, nil
        }

        os.Remove(outputPath)
    }

    // Tidak bisa capai target — return hasil terkecil yang bisa dicapai
    return domain.CompressResult{}, fmt.Errorf(
        "tidak bisa mencapai target %s, coba target yang lebih besar",
        formatMB(targetBytes),
    )
}

func formatMB(bytes int64) string {
    return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
