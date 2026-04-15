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
│   ├── domain/
│   │   ├── entity/
│   │   ├── repository/
│   │   └── service/
│   ├── usecase/
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── repository/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── websocket/
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
