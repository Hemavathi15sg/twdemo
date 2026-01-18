# Jira Integration Example

This document demonstrates how to use the Jira integration feature to fetch issue details.

## Setup

### 1. Environment Variables

Set the following environment variables:

```bash
export JIRA_BASE_URL=https://your-domain.atlassian.net
export JIRA_EMAIL=your-email@example.com
export JIRA_API_TOKEN=your-api-token
```

### 2. Generate Jira API Token

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Give it a name (e.g., "Grade Management API")
4. Copy the token

## Usage Examples

### Fetch Issue TEC-16

```bash
curl http://localhost:8080/api/jira/issues/TEC-16
```

**Successful Response (200 OK):**
```json
{
  "key": "TEC-16",
  "fields": {
    "summary": "Create enrollment feature",
    "description": "Implement student enrollment API with CRUD operations",
    "status": {
      "name": "Done"
    },
    "issuetype": {
      "name": "Story"
    },
    "created": "2026-01-15T10:00:00Z",
    "updated": "2026-01-18T07:00:00Z"
  }
}
```

### Error Responses

**Missing Configuration (500 Internal Server Error):**
```json
{
  "error": "JIRA_BASE_URL environment variable not set"
}
```

**Authentication Failed (401 Unauthorized):**
```json
{
  "error": "authentication failed: check JIRA_EMAIL and JIRA_API_TOKEN"
}
```

**Issue Not Found (404 Not Found):**
```json
{
  "error": "issue TEC-16 not found"
}
```

## Integration with Enrollment Feature

The Jira integration can be used to track requirements and stories related to the enrollment feature:

1. **TEC-16**: Create enrollment feature - This story tracks the implementation of the enrollment API
2. Use the Jira API to fetch issue details and display them in your application
3. Link enrollments to specific Jira issues for traceability

## Testing

Run the integration tests:

```bash
go test -tags=integration -v ./jira/...
```

## Implementation Details

- **Client**: Located in `jira/jira_client.go`
- **Handler**: Located in `handlers/jira_handler.go`
- **Endpoint**: `GET /api/jira/issues/{key}`
- **Authentication**: Uses Basic Auth with email and API token
- **Timeout**: 10 seconds
- **API Version**: Jira REST API v3
