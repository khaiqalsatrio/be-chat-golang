package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileUploadService struct {
	uploadDir string
}

func NewFileUploadService(uploadDir string) *FileUploadService {
	// Create upload directory if not exists
	os.MkdirAll(uploadDir, 0755)
	return &FileUploadService{
		uploadDir: uploadDir,
	}
}

type UploadResult struct {
	FilePath string `json:"file_path"`
	FileURL  string `json:"file_url"`
}

func (s *FileUploadService) UploadFile(file io.Reader, originalFilename string) (*UploadResult, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(originalFilename))
	validExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, // images
		".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".webm": true, // videos
	}

	if !validExts[ext] {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	// Generate unique filename
	filename := uuid.New().String() + ext

	// Create full file path
	filePath := filepath.Join(s.uploadDir, filename)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %v", err)
	}

	// Return relative path and URL
	return &UploadResult{
		FilePath: filepath.Join("uploads", filename),
		FileURL:  fmt.Sprintf("/uploads/%s", filename),
	}, nil
}
