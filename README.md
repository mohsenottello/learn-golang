# Layout

    cmd/server/        main(): config → store → router → serve
    internal/config/   env vars → Config struct
    internal/store/    the only package that runs SQL; one file per model + migrations/
    internal/mail/     SMTP
    internal/httpapi/  routes, middleware, handlers (one file per resource)
    web/               React + TypeScript (Vite); embed.go embeds the build

# Adding a resource (e.g. "projects")

1. `make migration name=projects` → fill in internal/store/migrations/0000N_projects.sql
2. internal/store/projects.go → struct + List/Get/Create/Update/Delete
3. internal/httpapi/projects.go → handlers
4. internal/httpapi/server.go → add the routes
5. web/src/api.ts + a page

Migrations run automatically when the server starts.

# Run

    cp .env.example .env && cp .env.local.example .env.local
    cd web && npm install && cd ..
    go mod tidy
    make dev     # terminal 1: Go API on :8080
    make web     # terminal 2: React on :5173
    make up      # or: everything in Docker on :8080
