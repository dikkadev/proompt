# QUICK REFERENCE - PROOMPT PROJECT

## 🏗️ PROJECT STATUS: FOUNDATION COMPLETE ✅
**Repository layer + Git integration working and tested**
**Next Phase: API layer implementation**

## CORE ENTITIES

### Prompt
- id (UUID), title, content, type, use_case
- model_compatibility_tags (JSON array)
- temperature_suggestion, other_parameters (JSON)
- git_ref (commit hash or branch)

### Snippet  
- id (UUID), title, content, description
- Global scope, can access prompt variables
- Cannot contain other snippets (1 layer max)

### Note
- id (UUID), prompt_id, title, body
- Multiple notes per prompt allowed

### Tags
- prompt_tags: many-to-many with prompts
- snippet_tags: many-to-many with snippets

## VARIABLE SYSTEM
- Syntax: `{{variable_name}}` or `{{var:default}}`
- String-only, no types
- Prompt variables override snippet variables
- Warnings for missing vars (not errors)
- Always computed on-demand, never stored resolved

## GIT INTEGRATION ✅ IMPLEMENTED
- **Single repo** with orphan branches: `prompts/{uuid}`, `snippets/{uuid}`, `notes/{uuid}`
- content.json stores entity data on each branch
- **Atomic DB + git operations** with transaction rollback
- Auto-managed, no git exposure to user
- **Location**: ~/.proompt/git-repo/ (single repo, multiple branches)

## KEY CONSTRAINTS
- No prompt execution
- No nested snippets
- Global snippets (not per-project)
- Use case is first-class field (not just tag)

## SEARCH FEATURES
- FTS5 on prompts, snippets, notes
- Tag-based filtering
- Use case filtering
- Future: semantic search for discovery

## FILE LOCATIONS ✅ IMPLEMENTED
- Database: ~/.proompt/database.db
- Git repo: ~/.proompt/git-repo/ (single repo)
- Branches: prompts/{uuid}, snippets/{uuid}, notes/{uuid}
- Each branch contains: content.json

## 🔧 BUILD & TEST
```bash
cd server/
make build                           # Build binary
go test ./internal/repository -v     # Test repository layer (✅ passing)
```

## PROMPT TYPES
- system, user, image, video (all same level)

## ORGANIZATION
- Primary: use_case field
- Secondary: arbitrary tags
- Tertiary: model compatibility tags
- Bidirectional prompt links (followup relationships)