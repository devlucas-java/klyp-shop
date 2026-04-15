
# 🧠 ARCHITECTURE FROM PROJECT

```md
# Architecture

This project follows Clean Architecture principles.

## 🧩 Layers

### 1. Domain
- Entities
- Interfaces (Repository contracts)
- Business rules

### 2. Usecase
- Application logic
- Orchestrates domain rules

### 3. Infrastructure
- Database (PostgreSQL + GORM)
- HTTP (Chi)
- WebSocket
- Security (JWT)

### 4. Delivery (Interface)
- REST API handlers
- Middleware

---

## 🔄 Flow

Request → Handler → Usecase → Domain → Repository → Database

---

## 🔐 Security

- JWT authentication
- Password hashing (bcrypt)
- WebSocket E2E encryption (planned)

---

## 📡 Future Improvements

- Microservice for Sellers
- Message Queue (RabbitMQ / Kafka)
- BitcoinPayment Gateway integration