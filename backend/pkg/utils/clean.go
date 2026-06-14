package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Panggil di awal main(), sebelum app.Listen
func CleanTmp(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return
	}
	for _, f := range files {
		os.Remove(f)
	}
	log.Println("tmp cleaned")
}

// Periodic cleanup setiap 1 jam — hapus file lebih dari 10 menit
func StartPeriodicCleanup(dir string) {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			files, _ := filepath.Glob(filepath.Join(dir, "*"))
			for _, f := range files {
				info, err := os.Stat(f)
				if err != nil {
					continue
				}
				if time.Since(info.ModTime()) > 10*time.Minute {
					os.Remove(f)
				}
			}
			log.Println("periodic tmp cleanup done")
		}
	}()
}
