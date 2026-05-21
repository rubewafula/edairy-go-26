package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveFile(file *multipart.FileHeader, folder string) (string, error) {
	if file == nil || file.Filename == "" || file.Size == 0 {
		return "", fmt.Errorf("invalid or empty file provided for folder %s", folder)
	}

	src, err := file.Open()
	if err != nil {
		log.Printf("SaveFile: Failed to open file from form (filename: %s): %v", file.Filename, err)
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	uploadDir := filepath.Join("uploads", folder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("SaveFile: Failed to create directory %s: %v", uploadDir, err)
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	name := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	path := filepath.Join(uploadDir, name)

	log.Printf("SaveFile: Attempting to create file: %s", path)

	dst, err := os.Create(path)
	if err != nil {
		log.Printf("SaveFile: Failed to create file on disk %s: %v", path, err)
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(path), nil
}
