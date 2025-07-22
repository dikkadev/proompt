# Keyboard-First UI Design for Proompt

## Core Principle
**Seamless keyboard + mouse interaction** - no modes, vim-inspired movement but always accessible via mouse too.

## Key Insight: Snippet Variable Inheritance
- Snippets can contain variables: `{{var_name}}`
- When snippet inserted into prompt, prompt variables take precedence
- Max 1 level deep (snippets can't contain other snippets)
- This creates a **variable dependency tree** that needs efficient navigation

## UI Architecture Concept

### Three-Panel Layout
```
[Sidebar]  [Main Editor]     [Preview/Inspector]
Library    Current Prompt    Live Output
Snippets   + Variables       + Dependencies
```

### Keyboard Navigation Philosophy
- **hjkl movement** between panels and within lists
- **Tab/Shift+Tab** for logical flow navigation
- **Enter** to activate/edit, **Escape** to cancel/back
- **Ctrl+** for global actions (save, search, etc.)
- **Alt+** for panel switching
- **No modal dialogs** - everything inline editable

## Specific Interactions

### Variable Editing
- **Click variable** → inline input field
- **Tab through variables** in logical order
- **Ctrl+Enter** to resolve/preview
- **Visual indicators** for provided/default/missing variables

### Snippet Insertion
- **Type @** → autocomplete dropdown
- **Arrow keys** to navigate suggestions
- **Enter** to insert
- **Drag from sidebar** also works
- **Show variable dependencies** immediately

### Panel Navigation
- **Alt+1/2/3** to focus panels
- **Ctrl+P** for command palette
- **Ctrl+F** for search within panel
- **Ctrl+Shift+F** for global search

### List Navigation (prompts, snippets, variables)
- **j/k** or **arrow keys** to move
- **Enter** to select/edit
- **Space** to preview
- **d** to delete (with confirmation)
- **t** to tag
- **l** to link

## Implementation Strategy

### 1. Keyboard Event System
- Global keyboard handler with context awareness
- Each component registers its shortcuts
- Prevent conflicts, show available shortcuts in status bar

### 2. Focus Management
- Clear visual focus indicators
- Logical tab order
- Focus trapping in modals/dropdowns
- Remember focus position when switching panels

### 3. Accessibility
- All keyboard shortcuts work with screen readers
- ARIA labels for dynamic content
- High contrast focus indicators
- Keyboard shortcuts documented in tooltips

## Technical Approach

### React Components
- **useKeyboard** hook for shortcut registration
- **FocusManager** context for panel coordination
- **CommandPalette** component for discoverability
- **InlineEdit** components for seamless editing

### State Management
- Current focus context
- Active shortcuts per context
- Undo/redo stack for all edits
- Auto-save with conflict resolution

This creates a **power user interface** that's discoverable for mouse users but blazingly fast for keyboard users.