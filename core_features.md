# Core Features - Completed Design

## What We're Building
A prompt management tool that sits between full execution platforms (Langfuse) and generic note apps. Focus on management without execution.

## Core Data Model
**Prompt Entity:**
- Content (the actual prompt text)
- Type (system, user, image, video - all same level)
- Use Case (first-class field for organization)
- Model compatibility tags
- Temperature/parameter suggestions
- Multiple notes (title + body pairs)
- Variables for templating

**Snippet Entity:**
- Reusable components that can be inserted into prompts
- Can access variables from the prompt they're inserted into
- Global scope (not per-project)
- Cannot contain other snippets (max 1 abstraction layer)

## Key Features

### 1. Templating System
- Variables in prompts: `{{variable_name}}`
- Snippet insertion with variable access
- One layer deep only (no nested snippets)
- Preview resolved output

### 2. Smart Metadata
- Configurable token count heuristics (user defines calculation methods)
- Model compatibility tracking
- Parameter suggestions (temperature, etc.)
- Flexible notes system (multiple title+body notes per prompt)

### 3. Organization
- Use case as primary organizational dimension
- Tags for additional categorization
- Search across content and metadata

### 4. Versioning (Future)
- Git-based versioning in background
- Auto-managed (rebase, etc.) - no git exposure to user
- Track changes over time

## Architecture Principles
- No prompt execution
- Simple but powerful composition
- Semantic organization over generic tagging
- Deep features that justify the tool's existence