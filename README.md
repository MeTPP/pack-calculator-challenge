# Pack Calculator

Calculates the optimal combination of packs to fulfill an order. Ships the fewest items possible using only whole packs, and among equal totals, uses the fewest packs.

Uses dynamic programming (coin change variant) with an optimization for large orders.

## Run

```sh
docker compose up --build
```

Open http://localhost:8080 or use the API:

```sh
# calculate packs
curl "http://localhost:8080/api/calculate?orderSize=501"
# [{"size":500,"count":1},{"size":250,"count":1}]

# view pack sizes
curl http://localhost:8080/api/pack-sizes

# update pack sizes
curl -X PUT -H "Content-Type: application/json" \
  -d '[23, 31, 53]' http://localhost:8080/api/pack-sizes

# edge case: order 500000 with packs [23, 31, 53]
curl "http://localhost:8080/api/calculate?orderSize=500000"
# [{"size":53,"count":9429},{"size":31,"count":7},{"size":23,"count":2}]
```

Without Docker:

```sh
go build -o pack-calculator ./cmd/main.go && ./pack-calculator
```

## API

| Method | Endpoint          | Description        |
|--------|-------------------|--------------------|
| GET    | /api/calculate    | Calculate packs    |
| GET    | /api/pack-sizes   | Get pack sizes     |
| PUT    | /api/pack-sizes   | Update pack sizes  |
| GET    | /health           | Health check       |

## Config

| Variable     | Default                  | Description          |
|-------------|--------------------------|----------------------|
| `PORT`      | `8080`                   | Server port          |
| `PACK_SIZES`| `250,500,1000,2000,5000` | Default pack sizes   |

## Test

```sh
go test -v -race ./...
```

## Structure

Clean architecture: domain -> usecases -> transport/repository.

```
cmd/main.go                    -- entry point
internal/
  domain/                      -- models, interfaces, errors
  usecases/                    -- business logic
  transport/http/              -- handlers, router, middleware
  repository/                  -- in-memory storage
  config/                      -- env config
templates/index.html           -- web UI
```
