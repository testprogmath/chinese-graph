#!/bin/bash
set -e

# This script sets up the VPS for first-time deployment
# Run as the 'deploy' user

echo "🚀 Setting up Chinese Learning Graph on VPS..."

# Create application directory
APP_DIR=~/apps/chinese-graph
mkdir -p $APP_DIR
cd $APP_DIR

# Clone repository
if [ ! -d ".git" ]; then
    echo "📦 Cloning repository..."
    git clone https://github.com/testprogmath/chinese-graph.git .
else
    echo "📦 Repository already exists, updating..."
    git fetch origin
    git reset --hard origin/main
fi

# Create .env file
if [ ! -f .env ]; then
    echo "🔧 Creating .env file..."
    cat > .env << 'EOF'
NEO4J_PASSWORD=chinesegraph123
NODE_ENV=production
JWT_SECRET=$(openssl rand -base64 32)
EOF
    echo "✅ .env file created (please update NEO4J_PASSWORD!)"
else
    echo "ℹ️ .env file already exists"
fi

# Create nginx directories
mkdir -p nginx/ssl

# Generate self-signed certificate for initial HTTPS (replace with Let's Encrypt later)
if [ ! -f nginx/ssl/cert.pem ]; then
    echo "🔒 Generating self-signed SSL certificate..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
    echo "✅ SSL certificate generated (replace with Let's Encrypt for production!)"
fi

# Pull Docker images
echo "🐳 Pulling Docker images..."
docker-compose -f docker-compose.prod.yml pull || docker compose -f docker-compose.prod.yml pull

# Start services
echo "🚀 Starting services..."
docker-compose -f docker-compose.prod.yml up -d --build || docker compose -f docker-compose.prod.yml up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check service health
echo "🔍 Checking service status..."
docker compose -f docker-compose.prod.yml ps

# Run initial data seeding
echo "🌱 Loading initial data..."
docker compose -f docker-compose.prod.yml exec -T backend sh -c "cd /app && ./server &"
sleep 5

# Show logs
echo "📋 Recent logs:"
docker compose -f docker-compose.prod.yml logs --tail=50

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update the .env file with a secure NEO4J_PASSWORD"
echo "2. Set up Let's Encrypt for proper SSL certificates:"
echo "   sudo certbot certonly --webroot -w /var/www/html -d yourdomain.com"
echo "3. Update nginx/conf.d/default.conf with your domain name"
echo "4. Configure GitHub secrets in your repository:"
echo "   - HOST: Your VPS IP address"
echo "   - USER: deploy"
echo "   - SSH_PRIVATE_KEY: Your SSH private key"
echo "   - KNOWN_HOSTS: Output of 'ssh-keyscan your-vps-ip'"
echo ""
echo "Access your application at:"
echo "   http://$(curl -s ifconfig.me)"
echo "   GraphQL Playground: http://$(curl -s ifconfig.me)/playground"