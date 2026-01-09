package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Config holds the configuration for MCP client
type Config struct {
	BaseURL        string
	APIKey         string
	Timeout        time.Duration
	MaxRetries     int
	EnableLogging  bool
}

// Client represents an MCP client
type Client struct {
	config     Config
	httpClient *http.Client
}

// EnrollmentRequest represents the enrollment data sent to MCP
type EnrollmentRequest struct {
	StudentID      int       `json:"student_id"`
	CourseID       int       `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
}

// EnrollmentResponse represents the response from MCP
type EnrollmentResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Timestamp string `json:"timestamp"`
}

// NewClient creates a new MCP client with the given configuration
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SendEnrollment sends enrollment data to MCP
func (c *Client) SendEnrollment(req EnrollmentRequest) (*EnrollmentResponse, error) {
	if c.config.EnableLogging {
		log.Printf("[MCP] Sending enrollment: StudentID=%d, CourseID=%d, Status=%s",
			req.StudentID, req.CourseID, req.Status)
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enrollment request: %w", err)
	}

	url := fmt.Sprintf("%s/api/enrollments", c.config.BaseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	httpReq.Header.Set("X-Client-Version", "1.0.0")

	var resp *http.Response
	var lastErr error

	// Retry logic
	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		if c.config.EnableLogging {
			log.Printf("[MCP] Attempt %d/%d to send enrollment", attempt, c.config.MaxRetries)
		}

		resp, lastErr = c.httpClient.Do(httpReq)
		if lastErr == nil && resp.StatusCode < 500 {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt < c.config.MaxRetries {
			backoff := time.Duration(attempt) * time.Second
			if c.config.EnableLogging {
				log.Printf("[MCP] Request failed, retrying after %v", backoff)
			}
			time.Sleep(backoff)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to send request after %d attempts: %w", c.config.MaxRetries, lastErr)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if c.config.EnableLogging {
			log.Printf("[MCP] Error response (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("MCP returned error status %d: %s", resp.StatusCode, string(body))
	}

	var enrollmentResp EnrollmentResponse
	if err := json.Unmarshal(body, &enrollmentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if c.config.EnableLogging {
		log.Printf("[MCP] Successfully sent enrollment: MessageID=%s", enrollmentResp.MessageID)
	}

	return &enrollmentResp, nil
}

// UpdateEnrollmentStatus sends enrollment status update to MCP
func (c *Client) UpdateEnrollmentStatus(studentID, courseID int, status string) error {
	if c.config.EnableLogging {
		log.Printf("[MCP] Updating enrollment status: StudentID=%d, CourseID=%d, Status=%s",
			studentID, courseID, status)
	}

	payload := map[string]interface{}{
		"student_id": studentID,
		"course_id":  courseID,
		"status":     status,
		"updated_at": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal status update: %w", err)
	}

	url := fmt.Sprintf("%s/api/enrollments/status", c.config.BaseURL)
	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send status update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if c.config.EnableLogging {
			log.Printf("[MCP] Status update failed (status %d): %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("MCP returned error status %d", resp.StatusCode)
	}

	if c.config.EnableLogging {
		log.Printf("[MCP] Successfully updated enrollment status")
	}

	return nil
}

// HealthCheck verifies connectivity to MCP
func (c *Client) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.config.BaseURL)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	if c.config.EnableLogging {
		log.Printf("[MCP] Health check successful")
	}

	return nil
}
