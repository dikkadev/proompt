# PROOMPT - QUICK REFERENCE

## 🚀 IMMEDIATE COMMANDS

### Start the Server
```bash
cd server
make run
# OR with Docker
make docker-run
```

### Run Tests
```bash
make test                    # Unit tests
make test-integration        # Docker integration tests
```

### Build
```bash
make build                   # Local binary
make docker-build           # Docker image
```

## 📡 API ENDPOINTS

**Base URL**: `http://localhost:8080`

### Prompts
- `GET /api/prompts` - List (supports ?type=, ?use_case=, ?limit=, ?offset=)
- `POST /api/prompts` - Create
- `GET /api/prompts/{id}` - Get by ID
- `PUT /api/prompts/{id}` - Update
- `DELETE /api/prompts/{id}` - Delete

### Snippets  
- `GET /api/snippets` - List
- `POST /api/snippets` - Create
- `GET /api/snippets/{id}` - Get by ID
- `PUT /api/snippets/{id}` - Update
- `DELETE /api/snippets/{id}` - Delete

### Notes
- `GET /api/prompts/{id}/notes` - List notes for prompt
- `POST /api/prompts/{id}/notes` - Create note
- `GET /api/notes/{id}` - Get note
- `PUT /api/notes/{id}` - Update note
- `DELETE /api/notes/{id}` - Delete note

## 📝 REQUEST EXAMPLES

### Create Prompt
```json
POST /api/prompts
{
  "title": "Code Review Assistant",
  "content": "Review this code for {{language}} and suggest improvements:\n\n{{code}}",
  "type": "system",
  "use_case": "code-review",
  "model_compatibility_tags": ["gpt-4", "claude-3"],
  "temperature_suggestion": 0.3
}
```

### Create Snippet
```json
POST /api/snippets
{
  "title": "Error Handling",
  "content": "Handle errors gracefully: {{error_context}}",
  "description": "Standard error handling pattern"
}
```

## 🔧 CONFIGURATION

Config file: `server/proompt.xml`
```xml
<config>
  <server host="localhost" port="8080"/>
  <database type="local" path="./data/proompt.db" migrations="./internal/db/migrations"/>
  <storage repos_dir="./data/repos"/>
</config>
```

## 🐳 DOCKER

### Development
```bash
docker compose -f server/compose.yaml up
```

### Testing
```bash
docker compose -f tests/docker/compose.test.yaml up --build
```

## 📁 KEY FILES

- `server/cmd/proompt/main.go` - Entry point
- `server/internal/api/server.go` - HTTP server setup
- `server/internal/api/handlers/` - API endpoints
- `server/internal/repository/` - Data layer
- `server/internal/models/` - Domain models
- `server/Dockerfile` - Container build
- `server/compose.yaml` - Orchestration