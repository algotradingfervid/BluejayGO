# Test Plan: Health Check Endpoint

## Summary
Verify health check endpoint returns correct JSON response with status and timestamp.

## Preconditions
- Server running on localhost:28090
- No authentication required

## User Journey Steps
1. Make request to GET /health
2. Verify response status code
3. Verify JSON structure
4. Verify response fields

## Test Cases

### Happy Path
- **Health endpoint responds**: Request GET /health, verify 200 status code
- **Content-Type header**: Verify Content-Type is application/json
- **JSON structure valid**: Verify response is valid JSON
- **Status field present**: Verify JSON contains "status" field
- **Status value correct**: Verify status field value is "ok"
- **Time field present**: Verify JSON contains "time" field
- **Time format valid**: Verify time field is valid RFC3339 format (e.g., "2024-01-15T10:30:45Z")
- **Time is current**: Verify time field reflects current server time (within reasonable margin)

### Response Format
- **Exact JSON structure**: Verify response matches: `{"status":"ok","time":"2024-...RFC3339..."}`
- **Field types**: Verify status is string, time is string
- **No extra fields**: Verify no unexpected additional fields in response
- **Consistent response**: Make multiple requests, verify consistent format

### Edge Cases
- **Multiple rapid requests**: Make multiple requests quickly, verify all return 200 and valid JSON
- **Time precision**: Verify RFC3339 format includes timezone (Z or offset)
- **JSON parsing**: Verify response is parseable by standard JSON libraries
- **Response encoding**: Verify UTF-8 encoding

## Selectors & Elements
- HTTP Response:
  - Status code: 200
  - Content-Type: application/json
  - JSON body with fields:
    - `status`: string value "ok"
    - `time`: string value in RFC3339 format

## HTMX Interactions
- None (standard HTTP GET request)

## Dependencies
- Health check handler: GET /health
- JSON encoding
- RFC3339 timestamp formatting
- No database dependency (health check should work even if DB is down, or explicitly check DB if that's part of health check)
