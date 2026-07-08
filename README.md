# Ticket System Backend (Go + Gin + SQLite + JWT)

Beginner-friendly ticket system backend.

## Endpoints

All responses are JSON.

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`
- `POST /tickets` (JWT required)
- `GET /tickets` (JWT required)
- `GET /tickets/{id}` (JWT required)
- `PATCH /tickets/{id}/status` (JWT required)

Ticket status flow:
- `open` -> `in_progress` -> `closed`
- `closed` cannot go back

Ownership rules:
- Users can create/view/update only their own tickets (no admin).

## Requirements

- Go (for local dev)
- SQLite (used via `database/app.db`)
- Docker (for container builds)

## Environment Variables

Create/edit `.env` (already included in the project):

- `JWT_SECRET` (required for JWT signing)
- `PORT` (optional, defaults to `8080`)

## Run Locally

1. Install dependencies:
   - `go mod tidy`
2. Start the server:
   - `go run .`
3. Server runs on:
   - `http://localhost:8080`

The app creates SQLite tables automatically on startup.

## Example: Register

```bash
curl -X POST http://localhost:8080/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Alice\",\"email\":\"alice@example.com\",\"password\":\"password123\"}"
```

## Example: Login

```bash
curl -X POST http://localhost:8080/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"alice@example.com\",\"password\":\"password123\"}"
```

This returns:
- `token` (JWT)

## Example: Create Ticket

```bash
curl -X POST http://localhost:8080/tickets ^
  -H "Authorization: Bearer YOUR_JWT_TOKEN" ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Bug in app\",\"description\":\"Something is broken\"}"
```

## Example: Update Ticket Status

```bash
curl -X PATCH http://localhost:8080/tickets/1/status ^
  -H "Authorization: Bearer YOUR_JWT_TOKEN" ^
  -H "Content-Type: application/json" ^
  -d "{\"status\":\"in_progress\"}"
```

## Docker (Render/Railway)

Build and run with Docker:

```bash
docker build -t ticket-system .
docker run -p 8080:8080 -e JWT_SECRET="your_secret_here" ticket-system
```

Deployment notes:
- Set `JWT_SECRET` as an environment variable on Render/Railway.
- The SQLite file is stored in the container at `database/app.db`. For production durability,
  you can add volume storage, but the assignment keeps it simple.

