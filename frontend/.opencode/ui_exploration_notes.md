# Frontend exploration - beautiful but disconnected

## UI is well-designed
3-panel layout that actually makes sense:
- Left: Snippet library (collapsible, categorized)
- Center: Editor + Live preview (60/40 split)  
- Right: Variable management (320px fixed)

Dark/light theme with custom accent colors. Smooth animations. Professional feel.

## The React stack is solid
- React 18 + TypeScript + Vite (modern)
- TanStack Query for server state (good choice)
- Radix UI primitives (50+ components, shadcn/ui style)
- Tailwind CSS v4 (latest)
- React Hook Form + Zod validation

## Deeper inspection shows mixed state

### What's actually connected:
- Health check works! Real API call with status indicator
- `usePrompts()` and `useSnippets()` hooks exist and are used
- API client is well-structured with proper error handling
- TypeScript types match backend models (mostly)

### What's still disconnected:
- Template preview does local regex parsing instead of using the amazing backend API
- Snippet sidebar expects `category`, `tags`, `isFavorite` fields that don't exist
- Mock data still lurking in components as fallbacks
- Variable resolution happens in frontend instead of using backend's sophisticated system

## The disconnect
Backend has this incredible template system:
```go
// Can track variable status, insert snippets, handle defaults
POST /api/template/preview
```

Frontend does this amateur hour regex:
```typescript
// Basic pattern matching, no real resolution
const variables = content.match(/\{\{([^}]+)\}\}/g)
```

Powerful backend capability not being used by frontend.

## SnippetSidebar.tsx is the poster child for this problem
Line 239: `const { data: snippetsData } = useSnippets();` ✅ Good!
Line 256: `s.isFavorite` ❌ This field doesn't exist in backend
Line 44: `const mockSnippets: Snippet[] = [` ❌ Still has fallback mock data

## Assessment progression
1. UI looks good
2. TanStack Query is properly configured
3. Mock data still present in some places
4. Backend template system not being used
5. Issues are fixable with proper integration

## What this tells me about working on it
- The UI foundation is excellent, don't mess with the design
- Focus on replacing local logic with backend API calls
- The data model mismatches are the biggest blocker
- Once connected properly, this will be a really nice tool

## Overall state
High quality pieces that aren't fully connected. Frontend and backend are both well-built but need proper integration.