# Rate Limiter + API Key Service (Go)

A lightweight **API Key management and rate limiting service** built with **Go**, designed as a learning-focused yet production-inspired backend project.

This project demonstrates:
- Idiomatic Go project structure
- Concurrency-safe rate limiting
- Middleware-based authentication
- Clean separation of concerns
- Practical API testing using Postman and curl

Repository:  
👉 https://github.com/princeofverry/rate-limiter-go

---

## ✨ Features

- 🔐 API Key generation and revocation
- 🚦 Per-key rate limiting using **Token Bucket** algorithm
- 🧵 Concurrency-safe (mutex-protected)
- 🧱 Middleware-based authentication & rate limiting
- ⚡ Fast, dependency-free (in-memory)
- 🧪 Easily testable with Postman or curl
- 📁 Clean, idiomatic Go project layout

---

## 🧠 Architecture Overview

```
Client
  │
  │  X-API-Key
  ▼
Middleware
  ├── API Key Validation
  └── Rate Limiter (Token Bucket)
  │
  ▼
Handlers
```

### Rate Limiting Strategy
- Algorithm: **Token Bucket**
- Default limit: **60 requests per minute per API key**
- Refill rate: continuous refill (per second)
- Storage: in-memory (per-process)

---

## 📂 Project Structure

```
rate-limiter-go/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── apikey/
│   │   └── store.go
│   ├── ratelimit/
│   │   └── limiter.go
│   └── httpapi/
│       ├── router.go
│       ├── handlers.go
│       └── middleware.go
├── go.mod
└── README.md
```

---

## 🚀 Getting Started

### Requirements
- Go **1.22+** (uses new `net/http` routing patterns)
- Git

### Clone the repository
```bash
git clone https://github.com/princeofverry/rate-limiter-go.git
cd rate-limiter-go
```

### Run the server
```bash
go run ./cmd/api
```

Server will start on:
```
http://localhost:8080
```

---

## 📌 API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "ok": true
}
```

---

### Create API Key
```
POST /v1/keys
```

Response:
```json
{
  "api_key": "your-generated-api-key"
}
```

---

### Revoke API Key
```
DELETE /v1/keys/{api_key}
```

Response:
```json
{
  "revoked": true
}
```

---

### Protected Endpoint
```
GET /v1/ping
```

Headers:
```
X-API-Key: your-api-key
```

Response:
```json
{
  "message": "pong"
}
```

---

## 🧪 Testing the API

### Using Postman

1. Create a new **Environment**
   - `base_url` → `http://localhost:8080`
   - `api_key` → generated key

2. Add header to protected requests:
   ```
   X-API-Key: {{api_key}}
   ```

3. Use **Collection Runner**
   - Iterations: `70`
   - Delay: `0 ms`
   - Environment: selected

Expected result:
- Requests 1–60 → `200 OK`
- Requests 61+ → `429 Too Many Requests`

---

### Using curl (Guaranteed Rate Limit Test)

```bash
for i in {1..70}; do
  curl -s -o /dev/null -w "%{http_code}\n" \
    -H "X-API-Key: YOUR_API_KEY" \
    http://localhost:8080/v1/ping
done
```

You should see `429` after exceeding the limit.

---

## ⚙️ Configuration

Currently configured directly in code:

```go
limiter := ratelimit.New(60, 60)
```

- Capacity: 60 tokens
- Refill: 60 tokens per minute

---

## 🧩 Design Decisions

- **In-memory storage** chosen for simplicity and learning purposes
- **Token Bucket** allows burst traffic while enforcing average rate
- **Middleware approach** keeps handlers clean and reusable
- **No external dependencies** for easier understanding and debugging

---

## 🎯 Future Goals / Roadmap

### Short Term
- [x] Add unit tests for rate limiter
- [ ] Add remaining-token visibility endpoint
- [ ] Structured logging (zap / zerolog)
- [ ] Graceful shutdown handling

### Mid Term
- [ ] Hash API keys instead of storing plaintext
- [ ] Persistent storage (PostgreSQL)
- [ ] Redis-based distributed rate limiter
- [ ] Configurable limits per API key

### Long Term
- [ ] Admin dashboard
- [ ] Prometheus metrics endpoint
- [ ] Docker & Docker Compose support
- [ ] API Gateway mode (reverse proxy)

---

## 📚 Learning Outcomes

This project reinforces:
- Go concurrency & mutex usage
- Middleware patterns
- Clean architecture with `internal/`
- HTTP server fundamentals
- Real-world rate limiting strategies

---

## 🧑‍💻 Author

**Verry Kurniawan**  
GitHub: https://github.com/princeofverry

---

## 📄 License

MIT License – feel free to use, modify, and learn from this project.
