# ⚽ Play5 Tournament API

A lightweight, email-free tournament management backend built with **Go** and **Fiber**. Designed for small soccer communities that currently organize everything through WhatsApp groups and need something simpler.

Players register with a username and password, fill out their profile, and wait for admin approval. No email verification, no friction. Just soccer.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [API Overview](#api-overview)
- [Authentication Flow](#authentication-flow)
- [Data Models](#data-models)
- [Internationalization (i18n)](#internationalization-i18n)
- [Environment Variables](#environment-variables)
- [Database Schema](#database-schema)
- [Refine.dev Frontend](#refinedev-frontend)
- [License](#license)

---

## Features

| Feature | Description |
|---------|-------------|
| **Email-free registration** | Username + password only. No SMTP, no verification headaches. |
| **Admin approval flow** | New accounts start pending. Admins approve via API. |
| **Editable profile (pre-approval)** | Users can update name, nickname, photo, phone, and positions while waiting. |
| **Soccer positions** | 1 or 2 preferred positions: Goalkeeper, Defender, Midfielder, Forward, Any. |
| **Bilingual support** | English and Spanish via `Accept-Language` header and user preference. |
| **JWT authentication** | Stateless tokens with role and approval status baked in. |
| **Admin assignment endpoint** | Promote users to admin via protected route. |
| **Photo upload** | Optional avatar support (local or S3-compatible storage). |

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.22+ |
| Framework | [Fiber](https://github.com/gofiber/fiber) v2 |
| ORM | [GORM](https://gorm.io/) |
| Database | PostgreSQL 15+ |
| Auth | JWT (golang-jwt/jwt/v5) |
| Password hashing | bcrypt |
| Validation | go-playground/validator/v10 |
| Storage | Local filesystem or S3/MinIO |
| i18n | go-i18n / custom middleware |

---

## Getting Started

### Prerequisites

- Go 1.22 or later
- PostgreSQL 15 or later
- (Optional) Docker & Docker Compose

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/Play5.git
cd Play5

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go

# Start the server
go run cmd/api/main.go
```

### Docker (Quick Start)

```bash
docker-compose up -d
```

This spins up PostgreSQL and the API together.

---

## Project Structure

```
.
cmd/
├── api/                    # Main Fiber application entrypoint
│   └── main.go
internal/
├── config/                 # Environment configuration
│   └── config.go
├── database/               # GORM setup and connection
│   └── database.go
├── handlers/               # HTTP route handlers
│   ├── auth_handler.go
│   ├── user_handler.go
│   └── admin_handler.go
├── middleware/             # Fiber middleware
│   ├── auth.go
│   ├── admin.go
│   ├── approval.go
│   └── i18n.go
├── models/                 # GORM models and enums
│   └── user.go
├── repository/             # Data access layer
│   └── user_repository.go
├── services/               # Business logic
│   ├── auth_service.go
│   ├── user_service.go
│   └── storage_service.go
├── utils/                  # Helpers (password, jwt, validators)
│   ├── password.go
│   ├── jwt.go
│   └── validator.go
└── locales/                # Translation files
    ├── active.en.toml
    └── active.es.toml

uploads/                    # Local photo storage (gitignored)

.env.example
Dockerfile
docker-compose.yml
go.mod
go.sum
README.md
```

---

## API Overview

### Base URL

```
http://localhost:3000/api/v1
```

### Public Routes (No Auth)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/register` | Create a new player account |
| `POST` | `/login` | Authenticate and receive JWT |

### Authenticated Routes (Any User)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/me` | Get current user profile |
| `PATCH` | `/me` | Update own profile (allowed even if pending) |
| `POST` | `/me/photo` | Upload profile photo |

### Admin Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/users` | List all users (supports `?isApproved=false`) |
| `GET` | `/users/:id` | Get specific user details |
| `PATCH` | `/users/:id/approve` | Approve or reject a registration |
| `PATCH` | `/users/:id/role` | Assign admin role |
| `GET` | `/positions` | List soccer positions (translated) |

---

## Authentication Flow

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│   Register  │────>│  Fill Profile │────>│  Await Approval  │
│  (username) │     │  (name, pic,  │     │  (admin via API  │
│  (password) │     │  positions...)  │     │   or WhatsApp)   │
└─────────────┘     └──────────────┘     └─────────────────┘
                                                │
                                                ▼
                                        ┌───────────────┐
                                        │  Full Access  │
                                        │  (tournament  │
                                        │   features)   │
                                        └───────────────┘
```

### JWT Claims

```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "username": "karel10",
  "role": "player",
  "isApproved": false,
  "preferredLanguage": "es",
  "iat": 1714320000,
  "exp": 1714323600
}
```

The `isApproved` flag in the JWT lets the frontend show a pending banner immediately after login, without an extra roundtrip.

---

## Data Models

### User

```go
type User struct {
    ID                uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
    Username          string      `gorm:"uniqueIndex;size:30;not null" json:"username"`
    PasswordHash      string      `gorm:"not null" json:"-"`
    Role              Role        `gorm:"type:varchar(20);default:'player'" json:"role"`

    // Approval lifecycle
    IsApproved        bool        `gorm:"default:false" json:"isApproved"`
    ApprovedAt        *time.Time  `json:"approvedAt,omitempty"`
    ApprovedBy        *uuid.UUID  `json:"approvedBy,omitempty"`

    // Profile
    FullName          string      `gorm:"size:100" json:"fullName"`
    Nickname          string      `gorm:"size:50" json:"nickname"`
    Phone             string      `gorm:"size:20" json:"phone"`
    PhotoURL          *string     `json:"photoUrl,omitempty"`

    // Soccer specific
    Positions         []Position  `gorm:"type:position[]" json:"positions"`
    PreferredLanguage string      `gorm:"default:'es'" json:"preferredLanguage"`

    CreatedAt         time.Time   `json:"createdAt"`
    UpdatedAt         time.Time   `json:"updatedAt"`
}
```

### Enums

```go
type Role string
const (
    RolePlayer Role = "player"
    RoleAdmin  Role = "admin"
)

type Position string
const (
    Goalkeeper Position = "goalkeeper"
    Defender   Position = "defender"
    Midfielder Position = "midfielder"
    Forward    Position = "forward"
    Any        Position = "any"
)
```

---

## Internationalization (i18n)

The API supports **English** and **Spanish** via the `Accept-Language` header or a `?lang=` query parameter.

### How It Works

1. **Request**: Client sends `Accept-Language: es` or `?lang=es`
2. **Middleware**: Detects and stores language in Fiber context
3. **User Preference**: Logged-in users can override via `preferredLanguage` field
4. **Response**: Error messages and dynamic content are translated accordingly

### Example: Positions Endpoint

**Request:**
```bash
curl http://localhost:3000/api/v1/positions?lang=es
```

**Response:**
```json
[
  { "value": "goalkeeper", "label": "Portero" },
  { "value": "defender", "label": "Defensa" },
  { "value": "midfielder", "label": "Mediocampista" },
  { "value": "forward", "label": "Delantero" },
  { "value": "any", "label": "Cualquiera" }
]
```

### Error Response Format

```json
{
  "error": {
    "code": "USERNAME_TAKEN",
    "message": "El usuario ya existe"
  }
}
```

The `code` is always machine-readable English. The `message` is human-translated.

---

## Environment Variables

```bash
# Server
PORT=3000
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=Play5
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-super-secret-key-min-32-chars
JWT_EXPIRATION_HOURS=24

# Storage
STORAGE_TYPE=local          # or "s3"
STORAGE_LOCAL_PATH=./uploads
# S3_BUCKET=
# S3_REGION=
# S3_ACCESS_KEY=
# S3_SECRET_KEY=
# S3_ENDPOINT=                # For MinIO compatibility

# Optional: WhatsApp notifications
WHATSAPP_ENABLED=false
WHATSAPP_API_URL=           # e.g., WhatsApp Business API or CallMeBot
WHATSAPP_ADMIN_NUMBER=      # Number to notify on new registrations
```

---

## Database Schema

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('player', 'admin');
CREATE TYPE position AS ENUM ('goalkeeper', 'defender', 'midfielder', 'forward', 'any');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(30) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'player',

    is_approved BOOLEAN NOT NULL DEFAULT false,
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id),

    full_name VARCHAR(100),
    nickname VARCHAR(50),
    phone VARCHAR(20),
    photo_url TEXT,

    positions position[] CHECK (array_length(positions, 1) BETWEEN 1 AND 2),
    preferred_language VARCHAR(5) NOT NULL DEFAULT 'es',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_pending ON users(is_approved) WHERE is_approved = false;
CREATE INDEX idx_users_username ON users(username);
```

---

## Refine.dev Frontend

This API is designed to pair with a [Refine.dev](https://refine.dev) frontend.

### Key Integration Points

| Refine Concept | API Mapping |
|----------------|-------------|
| `authProvider` | `POST /login`, `GET /me` |
| `dataProvider` | Standard REST endpoints for `users` resource |
| `i18nProvider` | Switch `preferredLanguage` on login, use `?lang=` param |
| Custom pages | "Complete Profile" page using `PATCH /me` (pre-approval) |
| Custom buttons | "Approve Player" calling `PATCH /users/:id/approve` |

### Refine Auth Provider Snippet

```tsx
const authProvider: AuthProvider = {
  login: async ({ username, password }) => {
    const { data } = await axios.post("/login", { username, password });
    localStorage.setItem("token", data.token);
    localStorage.setItem("locale", data.user.preferredLanguage);
    return { success: true, redirectTo: "/" };
  },
  check: async () => {
    const token = localStorage.getItem("token");
    if (!token) return { authenticated: false };

    const { data } = await axios.get("/me", {
      headers: { Authorization: `Bearer ${token}` }
    });

    return {
      authenticated: true,
      redirectTo: data.isApproved ? "/" : "/complete-profile",
    };
  },
  // ...
};
```

---

## License

MIT License. Built for the love of the game.

---

## Contributing

This is a community project. If your local soccer group needs features like team generation, match scheduling, or stat tracking, open an issue or PR.

**Vamos a jugar.** / **Let's play.**
