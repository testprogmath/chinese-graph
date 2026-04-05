#!/bin/bash

# Initialize git repository
echo "📝 Initializing Git repository..."
git init

# Add all files
echo "📦 Adding files..."
git add .

# Create initial commit
echo "💾 Creating initial commit..."
git commit -m "Initial commit: Chinese Learning Graph

- Backend with Go + GraphQL (gqlgen)
- Neo4j for graph database
- Docker setup for local development
- CI/CD with GitHub Actions
- Project board automation"

# Add remote repository
echo "🔗 Adding remote repository..."
git remote add origin https://github.com/testprogmath/chinese-graph.git

# Create main branch
git branch -M main

# Push to GitHub
echo "🚀 Pushing to GitHub..."
echo "You need to create the repository on GitHub first!"
echo "Go to: https://github.com/new"
echo ""
echo "Repository name: chinese-graph"
echo "Description: Visual graph-based Chinese learning application"
echo "Public/Private: Your choice"
echo ""
echo "After creating the repository, run:"
echo "git push -u origin main"

# Instructions for creating the repo
cat << EOF

=== Manual Steps Required ===

1. Go to https://github.com/new
2. Create repository with:
   - Repository name: chinese-graph
   - Description: Visual graph-based Chinese learning application
   - Initialize WITHOUT README, .gitignore, or license

3. After creating, run:
   git push -u origin main

4. Set up GitHub Secrets (Settings > Secrets):
   - CODECOV_TOKEN (if using Codecov)
   - Any deployment secrets

5. Enable GitHub Actions:
   - Go to Actions tab
   - Enable workflows

6. Set up GitHub Projects:
   - Go to Projects tab
   - Create new project "Chinese Learning App"
   - Use Board template

EOF