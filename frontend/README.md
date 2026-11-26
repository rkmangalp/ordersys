# OrderSys Frontend — API Test UI

React + Vite frontend for testing the OrderSys Go backend endpoints.

## Quick Start

```bash
cd frontend
npm install
npm run dev
```

Opens at `http://localhost:5173` — Vite proxies `/api/*` to `http://localhost:8080` (Go server).

## Features

- **Health Check** — Tests `GET /health` endpoint
- **Products List** — Tests `GET /api/v1/products` endpoint
- **Request Logs** — Shows method, path, status, duration, timestamp
- **Server Status Indicator** — Green/red based on health check
- **Test All** — Runs both endpoints sequentially

## Project Structure

```
frontend/
├── package.json
├── vite.config.js      # Proxy config for /api -> localhost:8080
├── index.html
├── .gitignore
├── README.md
└── src/
    ├── main.jsx        # Entry point
    ├── App.jsx         # Main UI component
    └── index.css       # Styles
```

## Vite Proxy

In `vite.config.js`:
```js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

This means `fetch('/api/v1/products')` in the browser goes to `http://localhost:8080/api/v1/products` on the Go server.

## Requirements

- Go backend running on `http://localhost:8080`
- Node.js 18+