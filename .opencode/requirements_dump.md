# Proompt Frontend Requirements & Constraints Dump

## Core System Understanding

<backend_features>
- Variable resolution: `{{variable_name}}` and `{{variable:default_value}}`
- Snippet insertion: `@snippet_name` and `@{snippet with spaces}`
- Snippets can contain variables (max 1 level deep)
- Variable precedence: prompt variables override snippet defaults
- Bidirectional prompt linking system
- Tag management for organization
- Template processing API with live preview/analyze endpoints
- Git-based versioning with orphan branches
- Full CRUD operations for prompts/snippets
</backend_features>

<variable_inheritance>
- Snippets are things that can go into variable slots
- Snippets also have variables (all this goes max 1 layer deep)
- Variables are 'passed down' to snippets with precedence
- When snippet inserted into prompt, prompt variables take precedence
- Creates a variable dependency tree that needs efficient navigation
</variable_inheritance>

## UI/UX Requirements

<keyboard_first>
- Must be usable with keyboard only AND mouse
- Should seamlessly work together (no mode switching)
- Take inspiration from vim-like movement but NO MODES anywhere
- Not just commands but also movement and actions
- hjkl movement between panels and within lists
- Tab/Shift+Tab for logical flow navigation
- Enter to activate/edit, Escape to cancel/back
- Ctrl+ for global actions (save, search, etc.)
- Alt+ for panel switching
- No modal dialogs - everything inline editable
</keyboard_first>

<live_preview>
- Live preview with no key - just live
- Updates automatically as you type/edit
- Shows resolved content with variables filled in
- Shows variable dependency tree
- Visual indicators for provided/default/missing variables
</live_preview>

<user_workflows>
1. Quick prompt assembly - grab existing prompt, fill variables, maybe add snippets, use it
2. Snippet library management - build up reusable chunks, organize them, find them fast
3. Variable filling workflow - see what needs values, fill them efficiently
4. Not heavy prompt authoring - more composition than creation from scratch
</user_workflows>

## Technical Constraints

<frontend_architecture>
- Pure client-side React app (no SSR)
- Will eventually be embedded in Go binary
- Bun + Vite + React + TypeScript
- Tailwind CSS v4 (new @import "tailwindcss" syntax)
- Vite proxy for API calls (/api → localhost:8080)
- Frontend running on localhost:5173 with hot reload
</frontend_architecture>

<api_integration>
- POST /api/template/preview - Full resolution with variables
- POST /api/template/analyze - Analysis without variable resolution
- Response includes: Resolved content, variable status, warnings
- CRUD endpoints for prompts, snippets, tags, links
- Real-time template processing for live preview
</api_integration>

## UI Design Concepts (Explored)

<three_panel_layout>
```
[Sidebar]  [Main Editor]     [Preview/Inspector]
Library    Current Prompt    Live Output
Snippets   + Variables       + Dependencies
```
- Left sidebar: prompt library + snippet browser
- Main area: current prompt with variable inputs
- Right panel: live preview/output
- Alt+1/2/3 for panel switching
</three_panel_layout>

<interaction_patterns>
- Fast access to existing prompts (search/filter/browse)
- Efficient variable filling (form-like or inline editing)
- Snippet picker/browser (quick insert, maybe drag-drop)
- Desktop-focused (can use more screen real estate)
- Command palette for quick actions (Cmd+K style)
- Inline variable editing (click placeholder → input field)
- Snippet drag-and-drop from sidebar
- Template preview that updates as you type
- Quick actions bar (save, copy, clear, etc.)
</interaction_patterns>

<keyboard_navigation>
- hjkl or arrow keys to move within panels
- Tab through variables in logical order
- Type @ → autocomplete dropdown for snippets
- Arrow keys to navigate suggestions
- Enter to insert, Escape to cancel
- Space to preview without editing
- / → quick filter current list
- Ctrl+P → command palette
- Ctrl+F → search within panel
- Ctrl+Shift+F → global search
</keyboard_navigation>

## Implementation Details

<react_hooks>
- useKeyboard hook for shortcut registration
- useFocus hook for panel coordination
- FocusManager context for panel coordination
- CommandPalette component for discoverability
- InlineEdit components for seamless editing
</react_hooks>

<state_management>
- Current focus context
- Active shortcuts per context
- Undo/redo stack for all edits
- Auto-save with conflict resolution
- Variable dependency tracking
- Live template resolution state
</state_management>

<accessibility>
- All keyboard shortcuts work with screen readers
- ARIA labels for dynamic content
- High contrast focus indicators
- Keyboard shortcuts documented in tooltips
- Clear visual focus indicators
- Logical tab order
- Focus trapping in modals/dropdowns
- Remember focus position when switching panels
</accessibility>

## Current Status

<completed>
- ✅ Backend: Production-ready with comprehensive API
- ✅ Frontend: Basic React/Tailwind setup working
- ✅ Three-panel layout implemented
- ✅ Keyboard event system and focus management
- ✅ Live preview with mock data
- ✅ Panel switching with Alt+1/2/3
- ✅ Basic variable editing interface
</completed>

<rejected>
- ❌ Current UI design (user feedback: "I hate it")
- Need to pivot to different approach
</rejected>

<next_considerations>
- Different layout pattern (not three-panel?)
- Different color scheme/visual design
- Different interaction model
- Simpler or more complex interface
- Different information architecture
</next_considerations>

## Open Questions

<design_decisions>
- What specific aspects of current UI are problematic?
- Preferred layout pattern (split-pane, tabs, dashboard, other)?
- Color preferences (dark/light theme, accent colors)?
- Information density (minimal vs detailed)?
- Primary use case priority (quick assembly vs library management)?
</design_decisions>