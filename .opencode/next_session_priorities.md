# NEXT SESSION PRIORITIES

## 🎯 PRIMARY GOAL: API LAYER IMPLEMENTATION

### Phase 1: HTTP Server Foundation
1. **Router Setup**
   - Choose router library (gorilla/mux, chi, or stdlib)
   - Basic server configuration
   - Graceful shutdown handling

2. **Middleware Stack**
   - Request logging middleware
   - Error handling middleware  
   - CORS middleware (if web UI planned)
   - Recovery middleware

### Phase 2: REST Endpoints
1. **Prompt Endpoints**
   ```
   GET    /api/prompts              # List with filtering
   POST   /api/prompts              # Create new
   GET    /api/prompts/{id}         # Get by ID
   PUT    /api/prompts/{id}         # Update
   DELETE /api/prompts/{id}         # Delete
   ```

2. **Snippet Endpoints** (same pattern)
3. **Note Endpoints** (same pattern)

### Phase 3: Request/Response Models
1. **Input DTOs** with validation
2. **Output DTOs** with proper JSON tags
3. **Error response** standardization
4. **Pagination** for list endpoints

### Phase 4: Handler Implementation
1. **Connect to repository layer**
2. **Input validation**
3. **Error handling and status codes**
4. **Response formatting**

### Phase 5: Testing
1. **HTTP endpoint tests**
2. **Integration tests**
3. **Error scenario testing**

## 🔧 TECHNICAL DECISIONS NEEDED

1. **Router Library**: gorilla/mux vs chi vs stdlib
2. **JSON Library**: stdlib vs faster alternatives
3. **Validation**: extend current validator or HTTP-specific
4. **Error Format**: RFC 7807 Problem Details vs custom
5. **Pagination**: cursor-based vs offset-based

## 📋 CLEANUP TASKS

1. **Fix Go hints**: Replace `interface{}` with `any` in git service
2. **Add API documentation**: OpenAPI/Swagger spec
3. **Configuration**: Add HTTP server config options

## 🎯 SUCCESS CRITERIA

- [ ] HTTP server starts and accepts requests
- [ ] All CRUD endpoints working for all entities
- [ ] Proper error handling and status codes
- [ ] Integration tests passing
- [ ] API documentation available

## 🚀 STRETCH GOALS

- [ ] Variable resolution endpoint (`POST /api/prompts/{id}/resolve`)
- [ ] Search endpoints with FTS5 integration
- [ ] Health check endpoint
- [ ] Metrics endpoint

The foundation is complete - time to build the user-facing API!