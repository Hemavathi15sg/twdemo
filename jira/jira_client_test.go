//go:build integration
// +build integration

package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestJiraClient_GetIssue(t *testing.T) {
	// Test 1: Missing environment variables
	t.Run("Missing environment variables", func(t *testing.T) {
		// Save existing env vars
		oldBaseURL := os.Getenv("JIRA_BASE_URL")
		oldToken := os.Getenv("JIRA_API_TOKEN")
		oldEmail := os.Getenv("JIRA_EMAIL")
		
		// Clear env vars
		os.Unsetenv("JIRA_BASE_URL")
		os.Unsetenv("JIRA_API_TOKEN")
		os.Unsetenv("JIRA_EMAIL")
		
		defer func() {
			// Restore env vars
			if oldBaseURL != "" {
				os.Setenv("JIRA_BASE_URL", oldBaseURL)
			}
			if oldToken != "" {
				os.Setenv("JIRA_API_TOKEN", oldToken)
			}
			if oldEmail != "" {
				os.Setenv("JIRA_EMAIL", oldEmail)
			}
		}()
		
		client := NewJiraClient()
		_, err := client.GetIssue("TEC-16")
		
		if err == nil {
			t.Error("Expected error when environment variables are not set")
		}
		
		if err.Error() != "JIRA_BASE_URL environment variable not set" {
			t.Errorf("Expected 'JIRA_BASE_URL environment variable not set', got '%s'", err.Error())
		}
	})
	
	// Test 2: Mock Jira server with successful response
	t.Run("Successful response from Jira", func(t *testing.T) {
		// Create mock Jira server
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			if r.URL.Path != "/rest/api/3/issue/TEC-16" {
				t.Errorf("Unexpected path: %s", r.URL.Path)
			}
			
			// Return mock response
			response := JiraIssue{
				Key: "TEC-16",
				Fields: IssueFields{
					Summary:     "Create enrollment feature",
					Description: "Implement student enrollment API",
					Status: StatusInfo{
						Name: "Done",
					},
					IssueType: IssueType{
						Name: "Story",
					},
					Created: "2026-01-15T10:00:00.000+0000",
					Updated: "2026-01-18T07:00:00.000+0000",
				},
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer mockServer.Close()
		
		// Set environment variables
		os.Setenv("JIRA_BASE_URL", mockServer.URL)
		os.Setenv("JIRA_API_TOKEN", "test-token")
		os.Setenv("JIRA_EMAIL", "test@example.com")
		defer func() {
			os.Unsetenv("JIRA_BASE_URL")
			os.Unsetenv("JIRA_API_TOKEN")
			os.Unsetenv("JIRA_EMAIL")
		}()
		
		client := NewJiraClient()
		issue, err := client.GetIssue("TEC-16")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if issue.Key != "TEC-16" {
			t.Errorf("Expected key 'TEC-16', got '%s'", issue.Key)
		}
		
		if issue.Fields.Summary != "Create enrollment feature" {
			t.Errorf("Expected summary 'Create enrollment feature', got '%s'", issue.Fields.Summary)
		}
		
		if issue.Fields.Status.Name != "Done" {
			t.Errorf("Expected status 'Done', got '%s'", issue.Fields.Status.Name)
		}
	})
	
	// Test 3: Issue not found
	t.Run("Issue not found", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errorMessages":["Issue does not exist"]}`))
		}))
		defer mockServer.Close()
		
		os.Setenv("JIRA_BASE_URL", mockServer.URL)
		os.Setenv("JIRA_API_TOKEN", "test-token")
		os.Setenv("JIRA_EMAIL", "test@example.com")
		defer func() {
			os.Unsetenv("JIRA_BASE_URL")
			os.Unsetenv("JIRA_API_TOKEN")
			os.Unsetenv("JIRA_EMAIL")
		}()
		
		client := NewJiraClient()
		_, err := client.GetIssue("NOTFOUND-1")
		
		if err == nil {
			t.Error("Expected error for non-existent issue")
		}
		
		if err.Error() != "issue NOTFOUND-1 not found" {
			t.Errorf("Expected 'issue NOTFOUND-1 not found', got '%s'", err.Error())
		}
	})
	
	// Test 4: Unauthorized
	t.Run("Unauthorized", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errorMessages":["Authentication failed"]}`))
		}))
		defer mockServer.Close()
		
		os.Setenv("JIRA_BASE_URL", mockServer.URL)
		os.Setenv("JIRA_API_TOKEN", "invalid-token")
		os.Setenv("JIRA_EMAIL", "test@example.com")
		defer func() {
			os.Unsetenv("JIRA_BASE_URL")
			os.Unsetenv("JIRA_API_TOKEN")
			os.Unsetenv("JIRA_EMAIL")
		}()
		
		client := NewJiraClient()
		_, err := client.GetIssue("TEC-16")
		
		if err == nil {
			t.Error("Expected error for unauthorized access")
		}
		
		if err.Error() != "authentication failed: check JIRA_EMAIL and JIRA_API_TOKEN" {
			t.Errorf("Expected authentication error, got '%s'", err.Error())
		}
	})
}
