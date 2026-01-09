# MCP Integration Documentation

## Overview

This document describes the integration between the Grade Management API and the MCP (Master Control Program) system for enrollment management as specified in JIRA task TEC-16.

## Architecture

The MCP integration provides seamless communication between the enrollment API and the MCP system:

- **Asynchronous Communication**: Enrollment operations trigger notifications to MCP without blocking the main API response
- **Resilient Design**: Failed MCP communications are logged but don't prevent local enrollment operations
- **Secure Authentication**: All MCP requests use Bearer token authentication
- **Automatic Retries**: Failed requests are automatically retried with exponential backoff

## Configuration

### Environment Variables

The MCP client is configured using environment variables. Copy `.env.example` to `.env` and configure:

```bash
# Required Configuration
MCP_BASE_URL=https://mcp.example.com      # MCP system endpoint
MCP_API_KEY=your-secret-api-key-here      # Authentication token

# Optional Configuration
MCP_TIMEOUT=30                             # Request timeout in seconds (default: 30)
MCP_MAX_RETRIES=3                          # Max retry attempts (default: 3)
MCP_ENABLE_LOGGING=true                    # Enable detailed logging (default: false)
```

### Obtaining Credentials

Contact your MCP administrator or DevOps team to obtain:
1. MCP Base URL for your environment (development, staging, production)
2. API Key with appropriate permissions for enrollment management

## Features

### 1. Enrollment Creation Sync

When a new enrollment is created via `POST /api/enrollments`, the system:
- Creates the enrollment in the local database
- Sends enrollment data to MCP asynchronously
- Returns success to the client immediately
- Logs any MCP communication failures

**MCP Endpoint**: `POST {MCP_BASE_URL}/api/enrollments`

**Payload**:
```json
{
  "student_id": 123,
  "course_id": 456,
  "enrollment_date": "2026-01-09T10:27:02Z",
  "status": "active"
}
```

### 2. Enrollment Status Updates

When enrollment status is updated via `PUT /api/enrollments/{id}`, the system:
- Updates the local enrollment record
- Notifies MCP of the status change
- Includes timestamp information

**MCP Endpoint**: `PUT {MCP_BASE_URL}/api/enrollments/status`

**Payload**:
```json
{
  "student_id": 123,
  "course_id": 456,
  "status": "completed",
  "updated_at": "2026-01-09T10:27:02Z"
}
```

### 3. Health Monitoring

The system performs a health check on startup to verify MCP connectivity.

**MCP Endpoint**: `GET {MCP_BASE_URL}/health`

## Error Handling

The integration implements robust error handling:

1. **Configuration Errors**: If MCP credentials are not provided, the application starts without MCP integration
2. **Connection Failures**: Health check failures at startup disable MCP integration with warning logs
3. **Runtime Errors**: Communication failures during operation are logged but don't block API operations
4. **Timeout Handling**: Requests timeout after the configured duration (default: 30 seconds)
5. **Retry Logic**: Failed requests are retried up to 3 times (configurable) with exponential backoff

## Logging and Audit Trail

When `MCP_ENABLE_LOGGING=true`, the system logs:

- Enrollment creation attempts
- Status update notifications
- Request retries and failures
- Response message IDs
- Health check results

Example log output:
```
[MCP] Sending enrollment: StudentID=123, CourseID=456, Status=active
[MCP] Attempt 1/3 to send enrollment
[MCP] Successfully sent enrollment: MessageID=msg-abc-123
```

## Testing

### Without MCP (Local Development)

Simply don't set the MCP environment variables. The application will run normally without MCP integration:

```bash
go run main.go
```

Expected output:
```
MCP integration disabled: MCP_BASE_URL environment variable is required
To enable MCP integration, set MCP_BASE_URL and MCP_API_KEY environment variables
```

### With MCP (Integration Testing)

1. Set up environment variables:
```bash
export MCP_BASE_URL=https://mcp-dev.example.com
export MCP_API_KEY=dev-api-key
export MCP_ENABLE_LOGGING=true
```

2. Start the application:
```bash
go run main.go
```

Expected output:
```
✓ MCP connection established successfully
🚀 Grade Management API starting on port :8080
```

3. Test enrollment creation:
```bash
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 123,
    "course_id": 456,
    "status": "active"
  }'
```

4. Check logs for MCP integration activity

### Mock MCP Server

For testing without a real MCP instance, you can use the provided mock server script (if available) or set up a simple HTTP server that responds to the expected endpoints.

## Deployment Considerations

### Environment-Specific Configuration

Use different MCP endpoints for each environment:

- **Development**: `MCP_BASE_URL=https://mcp-dev.example.com`
- **Staging**: `MCP_BASE_URL=https://mcp-staging.example.com`
- **Production**: `MCP_BASE_URL=https://mcp.example.com`

### Security Best Practices

1. **Never commit** `.env` files or API keys to version control
2. Use secure secret management systems (AWS Secrets Manager, HashiCorp Vault, etc.)
3. Rotate API keys regularly
4. Use HTTPS for all MCP communications in production
5. Monitor and alert on MCP authentication failures

### Monitoring and Alerting

Monitor these metrics:
- MCP request success/failure rate
- MCP response times
- Retry counts
- Authentication errors

Set up alerts for:
- Sustained MCP connection failures
- High retry rates
- Authentication errors

## Troubleshooting

### MCP Connection Failures

**Symptom**: `Warning: MCP health check failed`

**Solutions**:
1. Verify `MCP_BASE_URL` is correct and reachable
2. Check network connectivity to MCP
3. Verify API key is valid
4. Check MCP system status

### Authentication Errors

**Symptom**: `MCP returned error status 401` or `403`

**Solutions**:
1. Verify `MCP_API_KEY` is correct
2. Check API key permissions with MCP administrator
3. Verify API key hasn't expired

### Timeout Issues

**Symptom**: Frequent timeout errors in logs

**Solutions**:
1. Increase `MCP_TIMEOUT` value
2. Check MCP system performance
3. Verify network latency to MCP

## API Contract

### Expected MCP Response Format

For enrollment creation/updates, MCP should respond with:

```json
{
  "success": true,
  "message_id": "unique-message-identifier",
  "timestamp": "2026-01-09T10:27:02Z"
}
```

### HTTP Status Codes

- `200 OK`: Successful status update
- `201 Created`: Successful enrollment creation
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Invalid or missing API key
- `403 Forbidden`: Insufficient permissions
- `500 Internal Server Error`: MCP system error
- `503 Service Unavailable`: MCP temporarily unavailable

## Support

For issues related to:
- **MCP Configuration**: Contact DevOps team
- **API Integration**: Contact backend team
- **MCP System**: Contact MCP administrator

## Version History

- **v1.0.0** (2026-01-09): Initial MCP integration implementation (TEC-16)
