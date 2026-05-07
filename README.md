# Klyp Shop

Klyp Shop is an e-commerce platform inspired by Amazon, where users can buy products, review sellers, and sellers manage their own stores.

## IMPORTANT

- install C compile for development 
- install Docker 

## 🚀 Tech Stack

- Golang (Go)
- Chi (HTTP router)
- GORM (ORM)
- PostgreSQL
- JWT Authentication
- Bitcoin Payment
- WebSocket (real-time chat)
- Docker & Docker Compose
- Swagger (API documentation)

## 📦 Features

- User & Seller system
- Product management
- Reviews & ratings
- Seller dashboard
- Real-time chat (E2E encrypted)
- Authentication with JWT

## 🐳 Running with Docker

```bash
docker-compose up --build

````

# 📁 Estrutura do projeto

````bash
klyp-shop/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── dto/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── websocket/
│   ├── domain/
│   │   ├── entity/
│   │   ├── repository/
│   ├── application/
│   │   └── service/
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── repository/
│   │   └── security/
│   └── configs/
├── pkg/
│   ├── logger/
│   └── utils/
├── docs/ (swagger)
├── docker/
│   └── Dockerfile
├── docker-compose.yml
├── .env
├── README.md
├── ARCHITECTURE.md
├── HELP.md
└── TASKS.md

````

# .env demostration

WEB_SERVER_PORT=8080

DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=klyp_db
DB_USER=postgres
DB_PASSWORD=postgres

JWT_SECRET=super-secret-key-very-secure-123456
JWT_EXPIRE_IN=15
JWT_REFRESH_EXPIRE_IN=1440

# BTCPay Server — fill after creating store in BTCPay dashboard
BTCPAY_BASE_URL=http://localhost:14142
BTCPAY_STORE_ID=
BTCPAY_API_KEY=
BTCPAY_WEBHOOK_SECRET=
