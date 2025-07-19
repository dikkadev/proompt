# QUICK REFERENCE - PROOMPT PROJECT

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

## GIT INTEGRATION
- One repo per prompt in ~/.proompt/repos/prompt-{uuid}/
- content.json stores prompt data
- Database + git operations must be atomic
- Auto-managed, no git exposure to user

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

## FILE LOCATIONS
- Database: ~/.proompt/database.db
- Git repos: ~/.proompt/repos/prompt-{uuid}/
- Each repo contains: .git/ and content.json

## PROMPT TYPES
- system, user, image, video (all same level)

## ORGANIZATION
- Primary: use_case field
- Secondary: arbitrary tags
- Tertiary: model compatibility tags
- Bidirectional prompt links (followup relationships)