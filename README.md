# Indexer - Production Blockchain Dashboard

A high-performance Ethereum and Bitcoin indexer with a modern React dashboard.

## Tech Stack

- **Backend**: Go (Gin, GORM, PostgreSQL)
- **Frontend**: React (Vite, Tailwind v4, Framer Motion)
- **Database**: PostgreSQL

## Key Features

- **Real-time Indexing**: Automatically fetches and stores the latest blocks and transactions.
- **Unified Branding**: Standardized identifiers for consistent cross-chain data.
- **Interactive Dashboard**: Clickable details for every block and transaction.
- **Production HUD**: Live height indicators and update pause toggle.

## Setup Instructions

### 1. Environment Configuration

Create a `.env` file in the project root with the following contents:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=indexer_db
PORT=8989

ETH_RPC_URL=https://mainnet.infura.io/v3/your_key
BTC_RPC_URL=http://user:pass@your_node:8332
```

### 2. Run Backend

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8989/api/v1`.

### 3. Run Frontend

```bash
cd frontend
npm install
npm run dev
```

The dashboard will be available at `http://localhost:5173`.

## Deployment

For production deployment:

1. **Build the Docker image**:
   ```bash
   docker build -t indexer-backend:latest .
   ```
2. **Run the container** (adjust ports and env vars as needed):
   ```bash
   docker run -d -p 8989:8080 \
     --env-file .env \
     --name indexer indexer-backend:latest
   ```
3. **Frontend**: Build the frontend with `npm run build` and serve the `dist` folder via Nginx or any static server.
4. Use a process manager like Docker Compose or systemd to orchestrate the Go backend, frontend, PostgreSQL, and Redis services.
