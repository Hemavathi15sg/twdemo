package usecases

import (
	"encoding/json"
	"fmt"
	"os"
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
type TEC16Reader struct{}

// NewTEC16Reader creates a new TEC16 reader
func NewTEC16Reader() *TEC16Reader {
	return &TEC16Reader{}
}

// ReadFile reads and parses a TEC16 format file
func (r *TEC16Reader) ReadFile(filepath string) (*TEC16Data, error) {
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
