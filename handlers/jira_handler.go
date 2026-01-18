package handlers

import (
	"log"
	"net/http"

	"grademanagement-demo/jira"

	"github.com/gorilla/mux"
)

// JiraHandler handles HTTP requests for Jira integration
type JiraHandler struct {
	jiraClient *jira.JiraClient
}

// NewJiraHandler creates a new Jira handler instance
func NewJiraHandler(client *jira.JiraClient) *JiraHandler {
	return &JiraHandler{
		jiraClient: client,
	}
}

// GetJiraIssue handles GET /api/jira/issues/{key}
func (h *JiraHandler) GetJiraIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueKey := vars["key"]

	if issueKey == "" {
		respondError(w, "issue key is required", http.StatusBadRequest)
		return
	}

	issue, err := h.jiraClient.GetIssue(issueKey)
	if err != nil {
		log.Printf("Failed to fetch Jira issue %s: %v", issueKey, err)
		
		// Determine appropriate status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "issue "+issueKey+" not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "authentication failed: check JIRA_EMAIL and JIRA_API_TOKEN" {
			statusCode = http.StatusUnauthorized
		}
		
		respondError(w, err.Error(), statusCode)
		return
	}

	respondJSON(w, issue, http.StatusOK)
}
