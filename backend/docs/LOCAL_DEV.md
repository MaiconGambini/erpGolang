# Backend Local Development

Run local infrastructure from the repository root:

```bash
docker compose -f docker-compose.dev.yml up -d
```

Run the API:

```bash
cd backend
go run ./cmd/api
```
