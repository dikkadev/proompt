-- Drop FTS tables first
DROP TABLE IF EXISTS notes_fts;
DROP TABLE IF EXISTS snippets_fts;
DROP TABLE IF EXISTS prompts_fts;

-- Drop indexes
DROP INDEX IF EXISTS idx_snippet_tags_snippet;
DROP INDEX IF EXISTS idx_snippet_tags_tag;
DROP INDEX IF EXISTS idx_prompt_tags_prompt;
DROP INDEX IF EXISTS idx_prompt_tags_tag;
DROP INDEX IF EXISTS idx_prompts_updated;
DROP INDEX IF EXISTS idx_prompts_type;
DROP INDEX IF EXISTS idx_prompts_use_case;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS prompt_links;
DROP TABLE IF EXISTS snippet_tags;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS prompt_tags;
DROP TABLE IF EXISTS snippets;
DROP TABLE IF EXISTS prompts;