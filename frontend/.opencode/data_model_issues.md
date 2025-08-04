# Frontend Data Model Issues

## 🔧 Snippet Schema Mismatch

### Current Frontend Expectations
```typescript
// SnippetSidebar.tsx expects:
interface Snippet {
  id: string;
  title: string;
  content: string;
  description?: string;
  created_at: string;
  updated_at: string;
  git_ref?: string;
  
  // ❌ These don't exist in backend:
  category?: string;      // Used for grouping in sidebar
  tags?: string[];        // Used for filtering
  isFavorite?: boolean;   // Used for favorites section
}
```

### Actual Backend Response
```typescript
// From /api/snippets:
interface Snippet {
  id: string;
  title: string;
  content: string;
  description?: string;
  created_at: string;
  updated_at: string;
  git_ref?: string;
  // That's it - no category, tags, or isFavorite
}
```

### Impact on UI Components

**SnippetSidebar.tsx Issues**:
- Line 256: `s.isFavorite` - undefined, breaks favorites filtering
- Line 261: `snippet.category` - undefined, breaks category grouping  
- Line 270: `s.tags.includes(tag)` - undefined, breaks tag filtering

**Current Workarounds**:
- Component has fallback logic but still tries to access missing fields
- Mock data still present as backup

## 🎯 Resolution Options

### Option A: Add Fields to Backend
**Pros**: Frontend works as designed
**Cons**: Backend schema changes, migration needed

**Implementation**:
- Add `category`, `is_favorite` columns to snippets table
- Add snippet_tags table (already exists?)
- Update API responses

### Option B: Modify Frontend (Recommended)
**Pros**: No backend changes, simpler
**Cons**: Need to rework UI organization

**Implementation**:
- Remove category-based grouping
- Use simple list or search-based organization
- Remove favorites feature (or implement differently)
- Use backend tags table if it exists

### Option C: Hybrid Approach
**Pros**: Best of both worlds
**Cons**: More complex

**Implementation**:
- Use existing tags table for categories (tag = category)
- Add is_favorite as boolean field
- Map tags to categories in frontend

## 🔍 Next Steps

1. **Check Backend Tags Implementation**
   - Does snippet_tags table exist?
   - Are tag endpoints implemented?
   - Can we use tags as categories?

2. **Simplify Frontend First**
   - Remove category/favorite dependencies
   - Get basic snippet list working
   - Add organization features later

3. **Test Current API**
   - Verify snippet CRUD operations
   - Check what fields are actually returned
   - Confirm data persistence