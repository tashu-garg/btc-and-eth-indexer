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

Edit the `.env` file in the root directory:

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

1. Build the frontend: `npm run build`
2. Serve the `dist` folder via Nginx or the Go server.
3. Use a process manager like Docker or systemd for the Go binary.
# btc-and-eth-indexer
