#!/bin/bash

# Script to load seed data in production
# Run this after deployment to populate initial data

echo "🌱 Loading seed data into Neo4j..."

# Build a temporary container with the seed script
docker compose -f docker-compose.prod.yml exec -T backend sh << 'EOF'
cd /app
if [ -f scripts/seed_data.go ]; then
    echo "Running seed data script..."
    go run scripts/seed_data.go
else
    echo "Seed data script not found in container"
fi
EOF

echo "✅ Seed data loading complete!"