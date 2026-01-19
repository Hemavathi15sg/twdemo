package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TEC16Data represents the structure of a TEC16 format file
type TEC16Data struct {
	Format      string              `json:"format"`
	Version     string              `json:"version"`
	Enrollments []TEC16Enrollment   `json:"enrollments"`
}

// TEC16Enrollment represents an enrollment entry in TEC16 format
type TEC16Enrollment struct {
	StudentID      int64     `json:"student_id"`
	CourseID       int64     `json:"course_id"`
	Status         string    `json:"status"`
	EnrollmentDate time.Time `json:"enrollment_date"`
}

// TEC16Reader handles reading and parsing TEC16 format files
type TEC16Reader struct {
	allowedDir string
}

// NewTEC16Reader creates a new TEC16 reader
func NewTEC16Reader() *TEC16Reader {
	// Default to current working directory if not set
	allowedDir := os.Getenv("TEC16_DATA_DIR")
	if allowedDir == "" {
		var err error
		allowedDir, err = os.Getwd()
		if err != nil {
			// Fallback to current directory (relative) if we can't get absolute path
			allowedDir = "."
		}
	}
	return &TEC16Reader{
		allowedDir: allowedDir,
	}
}

// ReadFile reads and parses a TEC16 format file
func (r *TEC16Reader) ReadFile(filepath string) (*TEC16Data, error) {
	// Validate file path to prevent directory traversal
	if err := r.validateFilePath(filepath); err != nil {
		return nil, err
	}

	// Read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var tec16Data TEC16Data
	if err := json.Unmarshal(data, &tec16Data); err != nil {
		return nil, fmt.Errorf("failed to parse TEC16 data: %w", err)
	}

	// Validate format
	if tec16Data.Format != "tec16" {
		return nil, fmt.Errorf("invalid format: expected 'tec16', got '%s'", tec16Data.Format)
	}

	return &tec16Data, nil
}

// validateFilePath validates that the file path is within allowed directory
func (r *TEC16Reader) validateFilePath(path string) error {
	// Clean the path to resolve any .. or .
	cleanPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Get absolute allowed directory
	absAllowedDir, err := filepath.Abs(r.allowedDir)
	if err != nil {
		return fmt.Errorf("failed to resolve allowed directory: %w", err)
	}

	// Use filepath.Rel to properly check if cleanPath is within absAllowedDir
	// This handles cross-platform path separators correctly
	relPath, err := filepath.Rel(absAllowedDir, cleanPath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// If the relative path starts with "..", it's outside the allowed directory
	if strings.HasPrefix(relPath, "..") || filepath.IsAbs(relPath) {
		return fmt.Errorf("file path not allowed: must be within %s", absAllowedDir)
	}

	// Check if filename (not full path) starts with a dot (hidden file)
	filename := filepath.Base(cleanPath)
	if strings.HasPrefix(filename, ".") {
		return fmt.Errorf("access to hidden files is not allowed")
	}

	return nil
}

// ToCreateEnrollmentRequests converts TEC16 enrollments to CreateEnrollmentRequest objects
func (r *TEC16Reader) ToCreateEnrollmentRequests(data *TEC16Data) []*CreateEnrollmentRequest {
	requests := make([]*CreateEnrollmentRequest, len(data.Enrollments))
	for i, e := range data.Enrollments {
		requests[i] = &CreateEnrollmentRequest{
			StudentID:      e.StudentID,
			CourseID:       e.CourseID,
			Status:         e.Status,
			EnrollmentDate: e.EnrollmentDate,
		}
	}
	return requests
}
