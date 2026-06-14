package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	delivery "shucompress/internal/controller/http"
	"shucompress/internal/repository"
	"shucompress/internal/usecase"
	"shucompress/pkg/compressor"
	"shucompress/pkg/utils"
)

func main() {
	utils.CleanTmp("./tmp")
    utils.StartPeriodicCleanup("./tmp")

	// 1. Init repository
	fileRepo := repository.NewFileRepository("./tmp")

	// 2. Init semua compressor
	pdfCompressor := compressor.NewPDFCompressor()
	imageCompressor := compressor.NewImageCompressor()
	pptxCompressor := compressor.NewPPTXCompressor()

	if !compressor.GhostscriptAvailable() {
		log.Println("warning: Ghostscript not found; PDF and PPTX compression will be unavailable")
	}

	// 3. Init usecase — inject semua dependency
	compressUC := usecase.NewCompressUsecase(
		fileRepo,
		pdfCompressor,
		imageCompressor,
		pptxCompressor,
	)

	// 4. Init handler
	compressHandler := delivery.NewCompressHandler(compressUC)

	// 5. Init Fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "POST",
		AllowHeaders: "Content-Type",
	}))

	// 6. Routes
	delivery.RegisterRoutes(app, compressHandler)

	log.Println("shucompress running on :8085")
	log.Fatal(app.Listen(":8085"))
}
