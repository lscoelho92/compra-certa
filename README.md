# Compra Certa

Compra Certa is a web app for managing grocery purchases, products, and categories. It focuses on a fast workflow to register purchases with multiple items and track monthly totals.

## Technologies

- Frontend: Nuxt 3, Vue 3, TypeScript, Tailwind CSS
- Backend: Go (Gin), GORM, PostgreSQL
- Local infra: Docker Compose

## Running locally

### 1) Backend + database (Docker)

```bash
docker compose up -d --build
```

- API: http://localhost:8080
- Version endpoint: http://localhost:8080/api

### 2) Frontend (Nuxt)

```bash
cd apps/web
npm install
npm run dev
```

- App: http://localhost:3000

## Backend scripts (API version)

In the backend, the API version is injected at build time via ldflags.

```bash
cd apps/api
make build VERSION=v1.2.3
make run VERSION=v1.2.3
```

If you do not provide `VERSION`, the default is `v1`.

## Project structure

- apps/api: Go API
- apps/web: Nuxt frontend
- docker-compose.yml: local infra (API + Postgres)

## Notes

- The `/api` endpoint returns the current backend version.
- API routes are available without a prefix and with a version prefix (e.g. `/api/v1`).
