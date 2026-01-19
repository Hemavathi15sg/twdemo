# Grade Management API - Demo Session 1

## 🚀 Clean Starting Point for AI Delegation Demo

### Quick Setup
```bash
# Initialize Git repository
git init
git add .
git commit -m "Initial demo setup for Session 1"

# Create feature branch (IMPORTANT for Copilot PR targeting)
git checkout -b feature/session1-demo

# Start Redis
docker-compose up -d

# Test basic server
go run main.go
curl http://localhost:8080
```

### 📚 API Endpoints

#### Enrollment CRUD Operations
- `POST /api/enrollments` - Create a single enrollment
- `GET /api/enrollments` - List all enrollments
- `GET /api/enrollments/{id}` - Get enrollment by ID
- `PUT /api/enrollments/{id}` - Update enrollment
- `DELETE /api/enrollments/{id}` - Delete enrollment

#### TEC16 Import Feature
- `POST /api/enrollments/import/tec16` - Import enrollments from TEC16 format file

**Security:** File paths are restricted to the working directory (configurable via `TEC16_DATA_DIR` environment variable). Directory traversal attempts are blocked.

**TEC16 Format Example:**
```json
{
  "format": "tec16",
  "version": "1.0",
  "enrollments": [
    {
      "student_id": 1001,
      "course_id": 101,
      "status": "active",
      "enrollment_date": "2024-01-15T10:00:00Z"
    }
  ]
}
```

**Usage:**
```bash
# Import enrollments from TEC16 file
curl -X POST http://localhost:8080/api/enrollments/import/tec16 \
  -H "Content-Type: application/json" \
  -d '{"filepath": "/path/to/tec16_sample.json"}'
```

**Response:**
- `201 Created` - All records imported successfully
- `200 OK` - Some records imported, some failed (partial success)
- `400 Bad Request` - All records failed or invalid file

### 🎯 Session 1 AI Agent Plan

**Act 1: CRUD Boilerplate** → **Cloud Coding Agent**
- Complete enrollment CRUD API
- Validation and error handling
- Repository pattern implementation

**Act 2: Performance Caching** → **Local Agent** 
- Redis integration
- Cache-aside pattern
- Performance optimization

**Act 3: Quality Gates** → **Local Agent**
- OpenAPI specification
- Contract validation
- Integration testing

**Act 4: Documentation** → **Background Agent**
- Godoc comments
- README generation
- API documentation

### ⚙️ Copilot Agent Setup Tips

**For Cloud Coding Agent:**
- Ensure you're on feature branch before prompting
- Agent will create PR against current branch's upstream
- Use detailed Jira-style requirements

**For Local Agent:**
- Great for iterative improvements
- Faster for testing and validation
- Perfect for incremental changes

**For Background Agent:**
- Ideal for documentation tasks
- Non-blocking workflow
- Can work while you demo other features

---
**Ready for 30-minute AI delegation demo!** 🤖✨