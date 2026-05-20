# 🧠 Arquitetura do Projeto

Este projeto segue os princípios de **Domain-Driven Design (DDD)** combinados com **Clean Architecture**.

---

## 🧩 Camadas

### 1. Domain (`internal/domain/`)
O núcleo da aplicação. Não depende de nenhuma outra camada.

| Subpasta | Responsabilidade |
|---|---|
| `entity/` | Modelos de domínio com estado e comportamento (métodos de negócio) |
| `policy/` | Regras de negócio que envolvem múltiplas entidades ou contextos |
| `enums/` | Tipos enumerados do domínio (Role, OrderStatus, PaymentStatus) |
| `errors/` | Erros de domínio tipados com código HTTP e mensagem |

**Regra:** entidades contêm comportamento próprio (`ChangePassword`, `CancelPending`, `MarkAsPaid`). Policies contêm regras que cruzam entidades ou impõem limites de negócio (`CanChat`, `CanCreate`, `CanManage`).

### 2. Application (`internal/application/service/`)
Orquestra o domínio. Coordena repositórios, aplica policies e retorna DTOs.

- Não contém regras de negócio — delega para entidades e policies
- Depende das interfaces de repositório (nunca das implementações)
- Um service por agregado: `UserService`, `OrderService`, `ProductService`, etc.

### 3. Infrastructure (`internal/infrastructure/`)
Implementações concretas de tudo que é externo ao domínio.

| Subpasta | Responsabilidade |
|---|---|
| `database/` | Implementações GORM dos repositórios |
| `repository/` | Interfaces de repositório (contratos do domínio) |
| `security/jwt/` | Geração e validação de tokens JWT |
| `btcpay/` | Cliente HTTP para o BTCPay Server |

### 4. Delivery (`internal/delivery/`)
Ponto de entrada da aplicação. Traduz HTTP ↔ domínio.

| Subpasta | Responsabilidade |
|---|---|
| `http/handler/` | Handlers HTTP — decodificam request, chamam service, retornam response |
| `http/router/` | Registro de rotas por módulo |
| `http/middleware/` | Auth JWT, verificação de role, métricas |
| `http/dto/` | Structs de request/response organizados por domínio (`duser/`, `dorder/`, etc.) |
| `http/dto/mapper/` | Conversão entre entidades e DTOs de resposta |
| `socket/` | WebSocket para chat em tempo real (Hub, Client, Message) |

### 5. Module (`internal/module/`)
Composição root — monta cada módulo injetando dependências (repositório → service → handler → router).

---

## 🔄 Fluxo de uma Request

```
HTTP Request
    │
    ▼
Middleware (JWT Auth / Role)
    │
    ▼
Handler  ──── decodifica DTO de request
    │
    ▼
Service  ──── aplica Policy ──── valida regras de negócio
    │
    ▼
Entity   ──── executa comportamento de domínio
    │
    ▼
Repository Interface
    │
    ▼
Database Implementation (GORM + YugabyteDB)
    │
    ▼
HTTP Response (DTO)
```

---

## 📐 Onde cada regra fica

| Tipo de regra | Onde fica | Exemplo |
|---|---|---|
| Comportamento de uma entidade | `entity/` | `order.CancelPending()`, `user.ChangePassword()` |
| Regra que envolve múltiplas entidades | `policy/` | `ChatPolicy.CanChat(sender, receiver)` |
| Limite de negócio (quota, ownership) | `policy/` | `AddressPolicy.CanCreate(existing)` |
| Orquestração e persistência | `service/` | buscar, validar policy, salvar, retornar DTO |
| Validação de formato de entrada | `dto/` | `req.Validate()` |
| Controle de acesso por role | `middleware/` | `RoleMiddleware` |

---

## 🔐 Segurança

- Autenticação via JWT (HS256)
- Hash de senha com bcrypt
- Verificação de assinatura HMAC-SHA256 nos webhooks do BTCPay
- Middleware de role para endpoints restritos (SELLER, ADMIN)

---

## 💳 Pagamento Bitcoin

- Integração com **BTCPay Server** via API REST
- Fluxo: `CreateInvoice` → usuário paga → webhook `InvoiceSettled` → order marcada como `paid`
- Assinatura do webhook validada com HMAC-SHA256

---

## 💬 Chat em Tempo Real

- WebSocket via `gorilla/websocket`
- Hub centralizado gerencia conexões ativas
- Regra de negócio: admin conversa com qualquer um; seller conversa com buyer (e vice-versa); seller não conversa com seller; buyer não conversa com buyer

---

## 🗂 Módulos do Sistema

| Módulo | Endpoints base |
|---|---|
| Auth | `/api/v1/auth` |
| User | `/api/v1/user` |
| Seller | `/api/v1/seller` |
| Product | `/api/v1/product` |
| Order | `/api/v1/order` |
| Cart | `/api/v1/cart` |
| Payment | `/api/v1/payment` |
| Dashboard | `/api/v1/dashboard` |
| Chat | `/api/v1/chat` |
| Featured | `/api/v1/featured` |
| Address | `/api/v1/address` |

---

## 🚀 Melhorias Futuras

- Microserviço para Sellers
- Message Queue (RabbitMQ / Kafka)
- Refresh token
- Notificações em tempo real
- Seed de dados
- Documentação Swagger completa
