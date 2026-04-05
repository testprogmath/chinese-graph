# 🚀 VPS Deployment Guide

Complete guide for deploying Chinese Learning Graph on a VPS.

## 📋 VPS Requirements

- **Minimum**: 2GB RAM, 1 CPU, 20GB SSD
- **Recommended**: 4GB RAM, 2 CPU, 40GB SSD
- **OS**: Ubuntu 22.04 LTS (recommended) or Debian 11

## 🛠️ Initial VPS Setup

### 1. Create Non-Root User

```bash
# As root
adduser deploy
usermod -aG sudo deploy
su - deploy
```

### 2. Setup SSH Key

```bash
# On your local machine
ssh-copy-id deploy@your-vps-ip

# On VPS - disable password auth
sudo nano /etc/ssh/sshd_config
# Set: PasswordAuthentication no
sudo systemctl reload sshd
```

### 3. Configure Firewall

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 7474/tcp  # Neo4j browser (optional, for dev)
sudo ufw enable
```

## 📦 Install Dependencies

### Docker & Docker Compose

```bash
# Install Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
newgrp docker

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify
docker --version
docker-compose --version
```

### Go Runtime

```bash
# Download Go
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz

# Install
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### Nginx

```bash
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx
```

## 🚀 Deploy Application

### 1. Clone Repository

```bash
cd ~
git clone https://github.com/yourusername/chinese-graph.git
cd chinese-graph
```

### 2. Setup Environment

```bash
# Create .env file
cat > backend/.env << 'EOF'
# Production Environment
NODE_ENV=production
PORT=8080

# Neo4j
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=your-secure-password-here

# Redis
REDIS_URL=redis://localhost:6379

# Security
JWT_SECRET=your-jwt-secret-here
CORS_ORIGIN=https://yourdomain.com
EOF

# Set permissions
chmod 600 backend/.env
```

### 3. Production Docker Compose

```bash
# Create production compose file
cat > backend/docker-compose.prod.yml << 'EOF'
version: '3.8'

services:
  neo4j:
    image: neo4j:5-enterprise
    container_name: chinese-graph-neo4j
    restart: unless-stopped
    ports:
      - "7687:7687"  # Bolt
      # - "7474:7474"  # Browser (comment out in production)
    environment:
      - NEO4J_AUTH=neo4j/${NEO4J_PASSWORD}
      - NEO4J_EDITION=community
      - NEO4J_dbms_memory_pagecache_size=1G
      - NEO4J_dbms_memory_heap_initial__size=1G
      - NEO4J_dbms_memory_heap_max__size=1G
    volumes:
      - neo4j_data:/data
      - neo4j_logs:/logs
      - neo4j_conf:/conf
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    container_name: chinese-graph-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

volumes:
  neo4j_data:
  neo4j_logs:
  neo4j_conf:
  redis_data:
EOF
```

### 4. Build and Start Services

```bash
cd backend

# Start infrastructure
docker-compose -f docker-compose.prod.yml up -d

# Wait for Neo4j
echo "Waiting for Neo4j..."
sleep 30

# Build Go app
go build -o server cmd/server/main.go

# Load initial data
go run scripts/seed_data.go
```

## 🔄 Systemd Service

Create service for auto-start:

```bash
# Create service file
sudo cat > /etc/systemd/system/chinese-graph.service << 'EOF'
[Unit]
Description=Chinese Graph Backend
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy/chinese-graph/backend
ExecStart=/home/deploy/chinese-graph/backend/server
Restart=always
RestartSec=5
StandardOutput=append:/var/log/chinese-graph/app.log
StandardError=append:/var/log/chinese-graph/error.log

# Environment
Environment=NODE_ENV=production
EnvironmentFile=/home/deploy/chinese-graph/backend/.env

[Install]
WantedBy=multi-user.target
EOF

# Create log directory
sudo mkdir -p /var/log/chinese-graph
sudo chown deploy:deploy /var/log/chinese-graph

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable chinese-graph
sudo systemctl start chinese-graph
sudo systemctl status chinese-graph
```

## 🌐 Nginx Configuration

### 1. Basic HTTP Setup

```bash
# Create Nginx config
sudo cat > /etc/nginx/sites-available/chinese-graph << 'EOF'
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    # Force HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL will be configured by certbot
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # GraphQL endpoint
    location /graphql {
        proxy_pass http://localhost:8080/query;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # GraphQL Playground (disable in production)
    location /playground {
        # Comment out in production
        proxy_pass http://localhost:8080/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }

    # Frontend (when ready)
    location / {
        root /var/www/chinese-graph;
        try_files $uri $uri/ /index.html;
    }
}
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/chinese-graph /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 2. SSL Setup with Let's Encrypt

```bash
# Get SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal
sudo certbot renew --dry-run
```

## 📊 Monitoring

### 1. Setup Logging

```bash
# Create log rotation
sudo cat > /etc/logrotate.d/chinese-graph << 'EOF'
/var/log/chinese-graph/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 deploy deploy
    sharedscripts
    postrotate
        systemctl reload chinese-graph
    endscript
}
EOF
```

### 2. Basic Monitoring Script

```bash
# Create monitoring script
cat > ~/monitor.sh << 'EOF'
#!/bin/bash

# Check services
check_service() {
    if systemctl is-active --quiet $1; then
        echo "✅ $1 is running"
    else
        echo "❌ $1 is down!"
        # Restart service
        sudo systemctl restart $1
        # Send alert (configure your notification method)
        # curl -X POST https://hooks.slack.com/... 
    fi
}

# Check endpoints
check_endpoint() {
    if curl -f -s $1 > /dev/null; then
        echo "✅ $1 is responding"
    else
        echo "❌ $1 is not responding!"
    fi
}

echo "=== Service Status ==="
check_service chinese-graph
check_service nginx
check_service docker

echo -e "\n=== Docker Status ==="
docker-compose -f ~/chinese-graph/backend/docker-compose.prod.yml ps

echo -e "\n=== Endpoint Status ==="
check_endpoint "http://localhost:8080/query"
check_endpoint "http://localhost:7687"

echo -e "\n=== Resource Usage ==="
df -h | grep -E '^/dev/'
free -h
docker stats --no-stream
EOF

chmod +x ~/monitor.sh

# Add to crontab for regular checks
(crontab -l 2>/dev/null; echo "*/5 * * * * /home/deploy/monitor.sh >> /var/log/chinese-graph/monitor.log 2>&1") | crontab -
```

## 🔄 Update Process

```bash
# 1. Pull latest code
cd ~/chinese-graph
git pull origin main

# 2. Backup database
docker exec chinese-graph-neo4j neo4j-admin database dump --to-path=/data/backups/backup-$(date +%Y%m%d).dump neo4j

# 3. Build new version
cd backend
go build -o server.new cmd/server/main.go

# 4. Swap binaries
mv server server.old
mv server.new server

# 5. Restart service
sudo systemctl restart chinese-graph

# 6. Verify
curl http://localhost:8080/query
```

## 🔐 Security Hardening

### 1. Neo4j Security

```bash
# Change default password immediately
docker exec -it chinese-graph-neo4j cypher-shell -u neo4j -p neo4j
# ALTER CURRENT USER SET PASSWORD FROM 'neo4j' TO 'your-secure-password';
```

### 2. Fail2ban Setup

```bash
sudo apt install fail2ban

# Configure for Nginx
sudo cat > /etc/fail2ban/jail.local << 'EOF'
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[nginx-limit-req]
enabled = true
EOF

sudo systemctl restart fail2ban
```

### 3. Regular Updates

```bash
# Create update script
cat > ~/update-system.sh << 'EOF'
#!/bin/bash
sudo apt update
sudo apt upgrade -y
sudo apt autoremove -y
docker system prune -af
EOF

chmod +x ~/update-system.sh

# Weekly updates
(crontab -l 2>/dev/null; echo "0 2 * * 0 /home/deploy/update-system.sh >> /var/log/updates.log 2>&1") | crontab -
```

## 🎯 Performance Tuning

### Neo4j Optimization

```bash
# Edit Neo4j config
docker exec -it chinese-graph-neo4j bash
vi /var/lib/neo4j/conf/neo4j.conf

# Add/modify:
# dbms.memory.pagecache.size=2g
# dbms.memory.heap.initial_size=1g
# dbms.memory.heap.max_size=2g
```

### System Optimization

```bash
# Increase file limits
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf

# TCP optimization
sudo sysctl -w net.core.somaxconn=65535
sudo sysctl -w net.ipv4.tcp_max_syn_backlog=65535
```

## ✅ Validation Checklist

- [ ] Services running: `systemctl status chinese-graph`
- [ ] Docker containers up: `docker-compose ps`
- [ ] Nginx working: `curl -I https://yourdomain.com`
- [ ] GraphQL responding: `curl http://localhost:8080/query`
- [ ] Neo4j accessible: Check logs `docker logs chinese-graph-neo4j`
- [ ] SSL valid: `curl https://yourdomain.com`
- [ ] Firewall configured: `sudo ufw status`
- [ ] Backups working: Check backup directory
- [ ] Monitoring active: Check monitor logs

## 🆘 Troubleshooting

### Service won't start
```bash
journalctl -u chinese-graph -n 50
sudo systemctl status chinese-graph
```

### Neo4j issues
```bash
docker logs chinese-graph-neo4j
docker exec -it chinese-graph-neo4j neo4j-admin memrec
```

### Port conflicts
```bash
sudo netstat -tulpn | grep :8080
sudo lsof -i :8080
```

### Reset everything
```bash
sudo systemctl stop chinese-graph
docker-compose -f docker-compose.prod.yml down -v
# Start fresh
```

## 📞 Support

- GitHub Issues: https://github.com/yourusername/chinese-graph/issues
- Documentation: Check README.md and docs/
- Logs: `/var/log/chinese-graph/`