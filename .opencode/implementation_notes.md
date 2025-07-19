# IMPLEMENTATION STRATEGY NOTES

## TECH STACK RECOMMENDATIONS

### Backend Options (My Gut Feeling)
1. **Go** - Fits the git integration well (go-git library mentioned), good for CLI tools
2. **Python** - Great for rapid prototyping, excellent SQLite support
3. **Node.js** - If web UI is planned, full-stack JS makes sense

### Database Layer
- SQLite with FTS5 (already decided)
- Need ORM or query builder
- Atomic transactions critical for database+git sync

### Git Integration
- go-git library mentioned in docs (suggests Go preference)
- Need to handle atomic operations carefully
- File watching for external git changes?

## CORE IMPLEMENTATION PHASES

### Phase 1: Foundation
- [ ] Database schema implementation
- [ ] Basic CRUD operations for prompts
- [ ] Git repo creation/management
- [ ] Content.json format handling

### Phase 2: Templating
- [ ] Variable parsing (`{{var}}` and `{{var:default}}`)
- [ ] Snippet system
- [ ] Variable dependency resolution
- [ ] Preview generation

### Phase 3: Search & Organization
- [ ] FTS5 integration
- [ ] Tag management
- [ ] Use case filtering
- [ ] Prompt linking

### Phase 4: UI Layer
- [ ] CLI interface OR
- [ ] Web interface OR  
- [ ] Desktop app

## CRITICAL IMPLEMENTATION DETAILS

### Git Sync Strategy
```
On prompt create:
1. BEGIN TRANSACTION
2. INSERT into database
3. Create git repo
4. Write content.json
5. Git commit
6. COMMIT TRANSACTION
7. If any step fails, ROLLBACK everything
```

### Variable Resolution Algorithm
```
1. Parse prompt content for {{var}} patterns
2. Parse snippet content for {{var}} patterns  
3. Merge variable lists (prompt vars override snippet vars)
4. Apply defaults where no value provided
5. Generate warnings for missing vars
6. Return resolved content + warnings
```

### File Structure
```
~/.proompt/
  database.db
  repos/
    prompt-{uuid}/
      .git/
      content.json
```

## TESTING STRATEGY IDEAS
- Unit tests for variable parsing
- Integration tests for database+git sync
- End-to-end tests for full workflows
- Performance tests for large prompt collections

## POTENTIAL GOTCHAS
1. Git repo corruption handling
2. Concurrent access to same prompt
3. Large prompt content performance
4. Variable circular dependencies (shouldn't happen with 1-layer limit)
5. Database migration strategy

## UI/UX CONSIDERATIONS
- Variable dependency visualization (red/yellow/green system from docs)
- Snippet browser with variable preview
- Search interface design
- Prompt editing with live preview

## PERFORMANCE THOUGHTS
- FTS5 should handle search well
- Git operations might be slow for large repos
- Consider background git operations
- Lazy loading for prompt lists

## SECURITY CONSIDERATIONS
- Local tool, so minimal security concerns
- File permissions on ~/.proompt/
- No network exposure planned
- Git repo integrity checks

---

# Configuration Validation Implementation (COMPLETED)

## What Was Implemented

Successfully replaced the monolithic `validate()` function with a modern, tag-based validation system using **go-playground/validator v10.27.0** (latest version).

## Key Improvements

### 1. Declarative Validation Tags
- Moved validation rules from imperative code to declarative struct tags
- Rules are now visible directly on field definitions
- Easy to understand what each field requires

### 2. Custom Validator for Complex Relationships
- Created `database_exclusive` custom validator
- Handles the "exactly one database type" business rule
- Extensible pattern for future complex validations

### 3. User-Friendly Error Messages
- Custom error formatting converts technical validation errors to readable messages
- Maps struct field paths to config file paths (e.g., `Config.Server.Port` → `server.port`)
- Specific error messages for each validation type

### 4. Comprehensive Test Coverage
- Tests cover all validation scenarios
- Validates both positive and negative cases
- Ensures error messages are helpful

## Code Structure

### Struct Tags Used
```go
type Config struct {
    Database Database `validate:"required,database_exclusive"`
    Storage  Storage  `validate:"required"`
    Server   Server   `validate:"required"`
}

type Server struct {
    Host string `validate:"required"`
    Port int    `validate:"required,min=1,max=65535"`
}

type TursoDatabase struct {
    URL   string `validate:"required,url"`
    Token string `validate:"required"`
}
```

### Custom Validator
```go
func validateDatabaseExclusive(fl validator.FieldLevel) bool {
    db := fl.Field().Interface().(Database)
    localSet := db.Local != nil
    tursoSet := db.Turso != nil
    return localSet != tursoSet // XOR: exactly one must be set
}
```

## Benefits Achieved

1. **Maintainability**: Validation rules are declarative and close to field definitions
2. **Scalability**: Easy to add new validation rules as config grows
3. **Readability**: Clear what each field requires at a glance
4. **Extensibility**: Custom validators handle any business logic
5. **Error Quality**: User-friendly error messages with proper field paths
6. **Performance**: Leverages optimized validation library
7. **Standards**: Uses industry-standard validation approach

## Future Extensions

Adding new validation rules is now trivial:

```go
// New field with validation
NewField string `validate:"required,email"`

// New custom validator
configValidator.RegisterValidation("custom_rule", validateCustomRule)
```

The validation system is now ready to scale with the project's growth without becoming unwieldy.