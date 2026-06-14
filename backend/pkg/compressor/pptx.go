package compressor

import (
	"fmt"
	"os/exec"
	"shucompress/internal/domain"
)

type PPTXCompressor struct{}

func NewPPTXCompressor() *PPTXCompressor {
	return &PPTXCompressor{}
}

func (p *PPTXCompressor) Supports(fileType domain.FileType) bool {
	return fileType == domain.FileTypePPTX
}

func (p *PPTXCompressor) Compress(req domain.CompressRequest) (domain.CompressResult, error) {
	if !GhostscriptAvailable() {
		return domain.CompressResult{}, fmt.Errorf(
			"PPTX compression membutuhkan Ghostscript. Install Ghostscript atau gunakan kompresi image saja",
		)
	}

	// PPTX compress: pakai LibreOffice convert ke PDF dulu, lalu Ghostscript compress
	// Step 1: Convert PPTX → PDF via LibreOffice
	pdfPath := req.TempPath + "_converted.pdf"

	convertCmd := exec.Command(
		"soffice",
		"--headless",
		"--convert-to", "pdf",
		"--outdir", fmt.Sprintf("%s", req.TempPath[:len(req.TempPath)-len("/"+req.OriginalName)]),
		req.TempPath,
	)

	if output, err := convertCmd.CombinedOutput(); err != nil {
		return domain.CompressResult{}, fmt.Errorf("libreoffice error: %w\n%s — pastikan LibreOffice terinstall", err, output)
	}

	// Step 2: Compress PDF hasil convert
	outputPath := req.TempPath + "_compressed.pdf"
	gsCmd := ghostscriptBinary()

	cmd := exec.Command(
		gsCmd,
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/ebook",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		fmt.Sprintf("-sOutputFile=%s", outputPath),
		pdfPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return domain.CompressResult{}, fmt.Errorf("ghostscript error: %w\n%s", err, output)
	}

	return domain.CompressResult{
		OutputPath: outputPath,
		Filename:   "compressed_" + req.OriginalName + ".pdf",
	}, nil
}
