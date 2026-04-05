# 🇨🇳 Chinese Learning Graph

Visual graph-based Chinese learning application with mind-map style connections between words, compounds, and phrases.

## 🚀 Quick Start (Local Development)

### Prerequisites
- Docker & Docker Compose
- Go 1.25+
- Node.js 20+
- Make (optional but recommended)

### 1. Start Infrastructure

```bash
# Start Neo4j and Redis
cd backend
docker-compose up -d

# Wait for Neo4j to be ready (check at http://localhost:7474)
# Login: neo4j / chinesegraph123
```

### 2. Run Backend

```bash
# In backend directory
go mod download
go run cmd/server/main.go

# Or with custom Neo4j password
NEO4J_PASSWORD=chinesegraph123 go run cmd/server/main.go
```

Backend will be available at:
- GraphQL Playground: http://localhost:8080
- GraphQL Endpoint: http://localhost:8080/query

### 3. Load Sample Data

```bash
# In another terminal
cd backend
go run scripts/seed_data.go
```

### 4. Test GraphQL Queries

Open http://localhost:8080 and try:

```graphql
# Get a word with relationships
query GetWord {
  word(hanzi: "好") {
    id
    hanzi
    pinyin
    meanings
    compounds {
      word {
        hanzi
        pinyin
      }
      phrase
      meaning
    }
  }
}

# Get word graph
query GetGraph {
  wordGraph(hanzi: "好", depth: 2) {
    centerWord {
      hanzi
      pinyin
    }
    nodes {
      word {
        hanzi
        pinyin
        meanings
      }
      distance
    }
    edges {
      source
      target
      type
      label
    }
  }
}

# Search words
query Search {
  searchWords(query: "你", limit: 10) {
    hanzi
    pinyin
    meanings
  }
}
```

## 📦 Using Makefile

```bash
# Install all dependencies
make install

# Start everything
make start

# Stop services
make stop

# Load sample data
make seed

# Run tests
make test

# Clean everything
make clean
```

## 🏗️ Architecture

```
chinese-graph/
├── backend/          # Go + GraphQL API
│   ├── cmd/         # Entry points
│   ├── graph/       # GraphQL schema & resolvers
│   ├── internal/    # Business logic
│   └── scripts/     # Utils & seeders
├── frontend/        # React + TypeScript (coming soon)
└── mobile/         # React Native (future)
```

## 🔧 Configuration

### Environment Variables

```bash
# Backend
PORT=8080
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=chinesegraph123
REDIS_URL=localhost:6379

# Frontend (future)
REACT_APP_API_URL=http://localhost:8080/query
```

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐛 Troubleshooting

### Neo4j won't start
```bash
# Check logs
docker-compose logs neo4j

# Reset data (WARNING: deletes everything)
docker-compose down -v
docker-compose up -d
```

### Port already in use
```bash
# Change port
PORT=3001 go run cmd/server/main.go
```

### Can't connect to Neo4j
1. Wait 30 seconds after `docker-compose up`
2. Check http://localhost:7474 
3. Verify password matches docker-compose.yml

## 🚢 Production Deployment (VPS)

### Basic VPS Setup

```bash
# 1. Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# 2. Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 3. Install Go
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# 4. Clone and run
git clone https://github.com/yourusername/chinese-graph.git
cd chinese-graph/backend
docker-compose up -d
go build -o server cmd/server/main.go
./server
```

### Production Docker Compose

Create `docker-compose.prod.yml`:
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "80:8080"
    environment:
      - NEO4J_URI=bolt://neo4j:7687
      - NEO4J_PASSWORD=${NEO4J_PASSWORD}
    depends_on:
      - neo4j
      - redis

  neo4j:
    image: neo4j:5
    volumes:
      - neo4j_data:/data
    environment:
      - NEO4J_AUTH=neo4j/${NEO4J_PASSWORD}
    
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  neo4j_data:
  redis_data:
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name chinese-graph.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }
}
```

## 📚 Tech Stack

- **Backend**: Go, GraphQL (gqlgen), Neo4j, Redis
- **Frontend**: React, TypeScript, React Flow (coming soon)
- **Mobile**: React Native (planned)
- **Infrastructure**: Docker, GitHub Actions

## 🤝 Contributing

1. Create issue in GitHub
2. Fork the repo
3. Create feature branch
4. Make changes
5. Run tests
6. Submit PR

## 📄 License

MIT