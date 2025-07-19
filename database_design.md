# Database Design

```sql
-- Simplified schema for prompt management tool

-- Main prompts table
CREATE TABLE prompts (
    id TEXT PRIMARY KEY,  -- UUID
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('system', 'user', 'image', 'video')),
    use_case TEXT,  -- First-class field for organization
    model_compatibility_tags JSON,  -- Array of compatible models
    temperature_suggestion REAL,
    other_parameters JSON,  -- Other parameter suggestions
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Git reference (either commit hash or branch name depending on strategy)
    git_ref TEXT  -- e.g., "prompt-abc123" for branch or commit hash
);

-- Prompt tags (many-to-many)
CREATE TABLE prompt_tags (
    prompt_id TEXT REFERENCES prompts(id) ON DELETE CASCADE,
    tag_name TEXT,
    PRIMARY KEY (prompt_id, tag_name)
);

-- Notes attached to prompts (multiple title+body pairs)
CREATE TABLE notes (
    id TEXT PRIMARY KEY,
    prompt_id TEXT REFERENCES prompts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    body TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Reusable snippets (global scope)
CREATE TABLE snippets (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    git_ref TEXT
);

-- Snippet tags
CREATE TABLE snippet_tags (
    snippet_id TEXT REFERENCES snippets(id) ON DELETE CASCADE,
    tag_name TEXT,
    PRIMARY KEY (snippet_id, tag_name)
);

-- Bidirectional prompt links
CREATE TABLE prompt_links (
    from_prompt_id TEXT REFERENCES prompts(id) ON DELETE CASCADE,
    to_prompt_id TEXT REFERENCES prompts(id) ON DELETE CASCADE,
    link_type TEXT DEFAULT 'followup',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (from_prompt_id, to_prompt_id)
);

-- Performance indexes
CREATE INDEX idx_prompts_use_case ON prompts(use_case);
CREATE INDEX idx_prompts_type ON prompts(type);
CREATE INDEX idx_prompts_updated ON prompts(updated_at);

-- Tag search indexes
CREATE INDEX idx_prompt_tags_tag ON prompt_tags(tag_name);
CREATE INDEX idx_prompt_tags_prompt ON prompt_tags(prompt_id);
CREATE INDEX idx_snippet_tags_tag ON snippet_tags(tag_name);
CREATE INDEX idx_snippet_tags_snippet ON snippet_tags(snippet_id);

-- Full-text search virtual tables
CREATE VIRTUAL TABLE prompts_fts USING fts5(
    title, content, use_case, content=prompts, content_rowid=rowid
);

CREATE VIRTUAL TABLE snippets_fts USING fts5(
    title, content, description, content=snippets, content_rowid=rowid
);

CREATE VIRTUAL TABLE notes_fts USING fts5(
    title, body, content=notes, content_rowid=rowid
);
```