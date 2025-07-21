# Frontend Architecture Decisions

## Core Decisions Made
- **Web + TUI** - Web first, TUI later
- **Pure Client-Side React** - No server-side rendering, no Next.js server features
- **Static Bundle** - Builds to static files that can be served anywhere
- **Embedded Deployment** - Eventually embed static files in Go binary via embed
- **Makefile Integration** - Build process integrated into project root Makefile
- **Development Setup** - Server running separately, frontend with hot reload

## Technical Implications
- Use Create React App or Vite for pure client-side setup
- All API calls to existing Go backend (localhost:8080 in dev)
- Build output goes to server directory for embedding
- No SSR, no API routes in frontend
- CORS handling needed for dev (frontend port != backend port)

## Development Workflow
- Backend: `make run` or similar (localhost:8080)
- Frontend: `bun run dev` with hot reload (localhost:5174)
- Production: Static files embedded in Go binary, served by existing server

## Setup Complete ✅
- ✅ Vite + React + TypeScript project created
- ✅ Tailwind CSS v4 installed and configured
- ✅ Vite proxy configured for API calls (/api -> localhost:8080)
- ✅ Basic test page with Tailwind classes working
- ✅ Frontend running on localhost:5174

## Future Integration
- Go embed directive to include built frontend files
- Server serves static files at root or `/app` route
- Single binary deployment with both backend API and frontend UI