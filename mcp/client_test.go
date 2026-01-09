package mcp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := Config{
		BaseURL: "https://mcp.example.com",
		APIKey:  "test-api-key",
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout of 30s, got %v", client.config.Timeout)
	}

	if client.config.MaxRetries != 3 {
		t.Errorf("Expected default max retries of 3, got %d", client.config.MaxRetries)
	}
}

func TestSendEnrollment_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header with Bearer token")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json")
		}

		// Send successful response
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": true, "message_id": "msg-123", "timestamp": "2026-01-09T10:27:02Z"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := Config{
		BaseURL:       server.URL,
		APIKey:        "test-api-key",
		Timeout:       5 * time.Second,
		MaxRetries:    1,
		EnableLogging: false,
	}
	client := NewClient(config)

	// Test enrollment
	req := EnrollmentRequest{
		StudentID:      123,
		CourseID:       456,
		EnrollmentDate: time.Now(),
		Status:         "active",
	}

	resp, err := client.SendEnrollment(req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Error("Expected success to be true")
	}

	if resp.MessageID != "msg-123" {
		t.Errorf("Expected message_id 'msg-123', got '%s'", resp.MessageID)
	}
}

func TestSendEnrollment_Failure(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := Config{
		BaseURL:       server.URL,
		APIKey:        "test-api-key",
		Timeout:       5 * time.Second,
		MaxRetries:    1,
		EnableLogging: false,
	}
	client := NewClient(config)

	// Test enrollment
	req := EnrollmentRequest{
		StudentID:      123,
		CourseID:       456,
		EnrollmentDate: time.Now(),
		Status:         "active",
	}

	_, err := client.SendEnrollment(req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestHealthCheck_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected /health endpoint, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := Config{
		BaseURL:       server.URL,
		APIKey:        "test-api-key",
		Timeout:       5 * time.Second,
		EnableLogging: false,
	}
	client := NewClient(config)

	err := client.HealthCheck()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestLoadConfigFromEnv_Missing(t *testing.T) {
	// This test assumes environment variables are not set
	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Error("Expected error when MCP_BASE_URL is not set")
	}
}
