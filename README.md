# Compra Certa

Compra Certa is a web app for managing grocery purchases, products, and categories. It focuses on a fast workflow to register purchases with multiple items and track monthly totals.

## Tecnologias

- Frontend: Nuxt 3, Vue 3, TypeScript, Tailwind CSS
- Backend: Go (Gin), GORM, PostgreSQL
- Infra local: Docker Compose

## Como rodar localmente

### 1) Backend + banco (Docker)

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

## Scripts do backend (versao da API)

No backend, a versao da API e injetada no build via ldflags.

```bash
cd apps/api
make build VERSION=v1.2.3
make run VERSION=v1.2.3
```

Se nao informar `VERSION`, o padrao e `v1`.

## Estrutura do projeto

- apps/api: API em Go
- apps/web: frontend Nuxt
- docker-compose.yml: infraestrutura local (API + Postgres)

## Notas

- O endpoint `/api` retorna a versao atual do backend.
- As rotas da API estao disponiveis sem prefixo e com prefixo de versao (ex.: `/api/v1`).
