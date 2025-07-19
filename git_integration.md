# Git Integration Design

## Repository Structure

Each prompt gets its own git repository for versioning:
- One shadow repo per prompt (user's suggestion)
- Database stores current state and relationships
- Git repos store version history of prompt content

## Storage Layout

```
~/.proompt/
  database.db           # SQLite database
  repos/
    prompt-{uuid}/      # Individual git repos
      .git/
      content.json      # Prompt content + metadata
```

## Git Operations

### Creating a Prompt
1. Insert into database
2. Create new git repo in `repos/prompt-{uuid}/`
3. Write content.json with prompt data
4. Initial commit: "Create prompt: {title}"

### Updating a Prompt
1. Update database
2. Update content.json in prompt's repo
3. Commit: "Update prompt: {title}"

### Deleting a Prompt
1. Delete from database
2. Archive or delete git repo directory

## Content Format

Each prompt's git repo contains a single `content.json`:

```json
{
  "id": "uuid",
  "title": "Prompt Title",
  "content": "The actual prompt text...",
  "type": "system",
  "use_case": "code_review",
  "parameters": {
    "temperature": "0.3"
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Database-Git Sync

- Database is source of truth for current state
- Git repos are source of truth for history
- On startup: verify database matches latest git commits
- On change: update both database and git atomically

## Version History Access

- List commits in prompt's repo for history
- Checkout specific commits to view old versions
- Diff between commits to see changes
- No merge conflicts (single writer per repo)

## Implementation Notes

- Use go-git library for Git operations
- Atomic transactions: database + git together
- Background cleanup of old commits
- Optional: compress/archive old repos