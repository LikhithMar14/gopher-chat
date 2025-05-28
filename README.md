
# Gopher Chat 🐹💬

**Gopher Chat** is a blazing-fast, scalable, and secure chat server engineered in Go. Built for modern messaging demands, it features a clean API-first architecture and seamless developer experience. Perfect for teams, communities, and products seeking reliable real-time communication.

---

## 🚀 Key Features

- **API-First Design:** Clean, RESTful endpoints for all chat operations. Integrate with any frontend.
- **Modular Go Architecture:** Clear separation of API, business logic, and data layers.
- **Auto Database Migrations:** Effortless schema management on startup.
- **Centralized, Flexible Config:** Easy to deploy and customize.
- **API Key Security:** Protects sensitive endpoints out of the box.
- **OpenAPI (Swagger) Docs:** Self-describing, interactive API documentation.
- **Production-Ready:** Fast, robust, and built to scale.

---

## 🏗️ Architecture Overview

| Layer           | Description                                     |
|-----------------|-------------------------------------------------|
| **Client**      | Web/Mobile apps communicate via HTTP/REST       |
| **API Layer**   | Handles HTTP requests, validation, authentication |
| **Logic Layer** | Business rules, chat processing                 |
| **Storage**     | SQL database with automatic migrations          |
| **Config**      | Centralized, environment-driven configuration   |

**Data Flow Illustration:**

1. **Client** sends request (e.g., send message) to API.
2. **API Layer** authenticates and validates request.
3. **Logic Layer** processes the message.
4. **Storage Layer** persists the message in the database.
5. **API Layer** responds with success or error.

---

## 📦 Example: Sending a Message

**Request:**
```http
POST /v1/messages
Authorization: ApiKey <your-key>
Content-Type: application/json

{
  "room_id": "general",
  "text": "Hello, world!"
}
```

**Response:**
```json
{
  "id": "msg_12345",
  "room_id": "general",
  "user_id": "user_42",
  "text": "Hello, world!",
  "sent_at": "2025-05-28T15:00:00Z"
}
```

---

## 🚚 Quickstart

```bash
# 1. Clone the repo
git clone https://github.com/LikhithMar14/gopher-chat.git
cd gopher-chat

# 2. Configure environment (edit config.yaml or set ENV vars)
cp config.example.yaml config.yaml

# 3. Run database migrations automatically on startup

# 4. Launch the server
go run ./cmd/gopher-chat

# 5. Explore the API docs
# Visit: http://localhost:8080/swagger/
```

---

## 🔒 Security

- All sensitive endpoints are protected with API Key authentication.
- Never expose your API Keys publicly.
- Opt-in for HTTPS in production deployments.

---

## 🛠️ Advanced Configuration

- **Database:** PostgreSQL recommended (see config.yaml)
- **Logging:** Structured, environment-aware logs.
- **CORS:** Configurable for cross-origin support.

---

## 📝 API Documentation

- Auto-generated OpenAPI (Swagger) docs available at `/swagger/`.
- Explore, try, and integrate endpoints interactively.

---

## 🤝 Contributing

We welcome PRs and issues! To contribute:

1. Fork the repo.
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit and push your changes.
4. Open a pull request.

Please review our [CONTRIBUTING.md](CONTRIBUTING.md) and adhere to our code of conduct.

---

## 🧠 FAQ

- **Q: Can I run Gopher Chat in Docker?**
  - Yes! A production-ready Dockerfile is included.

- **Q: How do I add new endpoints?**
  - Implement a handler in the API layer and wire it up in the router.

- **Q: Is WebSocket/Realtime supported?**
  - Planned for a future release. Contributions welcome!

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.

---

## 🌐 Resources

- [OpenAPI Docs](http://localhost:8080/swagger/)
- [Go Documentation](https://golang.org/doc/)
- [GitHub Discussions](https://github.com/LikhithMar14/gopher-chat/discussions)

---

> **Built with passion, powered by Go. Happy chatting! 🚀**
