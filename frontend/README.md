# Proompt Frontend

A modern React application for prompt management and templating.

## Tech Stack

- **React 18** with TypeScript
- **Vite** for build tooling
- **Tailwind CSS** for styling
- **shadcn/ui** for component library
- **Bun** for package management
- **React Router** for navigation
- **React Query** for data fetching

## Development

### Prerequisites

- [Bun](https://bun.sh/) for package management

### Getting Started

```bash
# Install dependencies
bun install

# Start development server
bun run dev

# Build for production
bun run build

# Preview production build
bun run preview

# Run linting
bun run lint
```

## Features

- **Prompt Editor**: Rich text editing with variable and snippet support
- **Variable Management**: Dynamic variable detection and value input
- **Snippet Library**: Reusable content components with categorization
- **Live Preview**: Real-time template resolution and preview
- **Theme Support**: Dark/light mode with customizable accent colors

## Architecture

The frontend communicates with a Go-based API server for:
- Prompt CRUD operations
- Snippet management
- Template preview and analysis
- Variable resolution

## Project Structure

```
src/
├── components/          # React components
│   ├── ui/             # shadcn/ui components
│   ├── PromptEditor.tsx
│   ├── VariablePanel.tsx
│   └── SnippetSidebar.tsx
├── lib/                # Utilities
├── hooks/              # Custom React hooks
└── pages/              # Route components
```
