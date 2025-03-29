# Go File Share Platform 🚀

A secure, high-performance file-sharing platform built with Go, PostgreSQL, and Redis. Supports user authentication, file uploads, sharing links, and metadata search.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7+-DC382D?logo=redis)

## Features ✨

- 🔒 JWT Authentication (Register/Login)
- 📁 File Uploads (S3 or Local Storage)
- 🔗 Shareable Links with Expiration
- ⚡ Redis Caching for Metadata
- 🔍 Search Files by Name/Type
- 🧹 Background Cleanup Worker

## Tech Stack 🛠️

| Component      | Technology                |
|----------------|---------------------------|
| Backend        | Go (Gorilla Mux)          |
| Database       | PostgreSQL                |
| Cache          | Redis                     |
| Storage        | AWS S3 / Local Filesystem |
| Authentication | JWT                       |

## Prerequisites 📋

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- AWS CLI (for S3 - optional)

## Setup Guide 🛠️

### 1. Clone Repository
```bash
git clone https://github.com/fakubwoy/go-file-share.git
cd go-file-share
```

### 2. Configure Environment
Create .env file:
```bash
cp .env.example .env
```
Edit with your values:
```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=fileshare_user
DB_PASSWORD=securepassword
DB_NAME=fileshare
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your_secure_secret
```

### 3. Database Setup
```bash
sudo -u postgres psql -c "CREATE USER fileshare_user WITH PASSWORD 'securepassword';"
sudo -u postgres psql -c "CREATE DATABASE fileshare;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE fileshare TO fileshare_user;"
```

### 4. Run Migrations
```bash
sudo -u postgres psql -d fileshare -f migrations/001_init_schema.up.sql
```

### 5. Start Services
```bash
# Start Redis
sudo systemctl start redis-server

# Start PostgreSQL
sudo systemctl start postgresql
```

### 6. Run Application
```bash
go run cmd/main.go
```

## API Endpoints 🌐

| Method | Endpoint           | Description           |
|--------|--------------------|-----------------------|
| POST   | /register          | Register new user     |
| POST   | /login             | Login and get JWT token |
| POST   | /files             | Upload file           |
| GET    | /files             | List user's files     |
| POST   | /files/{id}/share  | Generate share link   |
| GET    | /share/{token}     | Access shared file    |

## Deployment Options 🚀

### Docker (Recommended)
```bash
docker-compose up --build
```

### AWS EC2
```bash
# Install dependencies
sudo apt update && sudo apt install -y golang postgresql redis

go build -o fileshare cmd/main.go
./fileshare
```