package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// JiraIssue represents a Jira issue
type JiraIssue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields represents the fields of a Jira issue
type IssueFields struct {
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Status      StatusInfo `json:"status"`
	IssueType   IssueType  `json:"issuetype"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
}

// StatusInfo represents the status of a Jira issue
type StatusInfo struct {
	Name string `json:"name"`
}

// IssueType represents the type of a Jira issue
type IssueType struct {
	Name string `json:"name"`
}

// JiraClient handles communication with Jira API
type JiraClient struct {
	BaseURL string
	Token   string
	Email   string
}

// NewJiraClient creates a new Jira client instance
func NewJiraClient() *JiraClient {
	return &JiraClient{
		BaseURL: os.Getenv("JIRA_BASE_URL"),
		Token:   os.Getenv("JIRA_API_TOKEN"),
		Email:   os.Getenv("JIRA_EMAIL"),
	}
}

// GetIssue retrieves a Jira issue by key
func (c *JiraClient) GetIssue(issueKey string) (*JiraIssue, error) {
	if c.BaseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL environment variable not set")
	}
	if c.Token == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN environment variable not set")
	}
	if c.Email == "" {
		return nil, fmt.Errorf("JIRA_EMAIL environment variable not set")
	}

	url := fmt.Sprintf("%s/rest/api/3/issue/%s", c.BaseURL, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Email, c.Token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue %s not found", issueKey)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed: check JIRA_EMAIL and JIRA_API_TOKEN")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var issue JiraIssue
	if err := json.Unmarshal(body, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &issue, nil
}
