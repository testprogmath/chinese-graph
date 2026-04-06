# Claude Code Development Notes

This document contains important gotchas, lessons learned, and best practices from the Chinese Graph project development and deployment process.

## Docker Build Gotchas

### 1. .dockerignore Location Must Match Build Context

**Problem**: Docker build failed with `stat /app/cmd/server/main.go: directory not found` even though the file existed locally.

**Root Cause**: The `.dockerignore` file was in the repository root (`./.dockerignore`) while the Docker build context was set to `./backend` in `docker-compose.prod.yml`.

**Solution**: Move `.dockerignore` to the build context directory (`./backend/.dockerignore`).

**Key Insight**: Docker looks for `.dockerignore` relative to the build context, not the repository root.

### 1.5. .dockerignore Whitelist Approach

**Problem**: Blacklist approach in .dockerignore can accidentally exclude needed directories.

**Solution**: Use whitelist approach with explicit directory patterns:
```dockerignore
# Ignore everything by default
*

# But include needed files for Go build
!go.mod
!go.sum
!gqlgen.yml

# Include directories and their contents
!cmd/
!cmd/**
!graph/
!graph/**
!internal/
!internal/**
```

**Why This Works Better**: 
- Explicit inclusion prevents accidental exclusions
- Easier to debug what's being copied
- More maintainable than long blacklists

```yaml
# docker-compose.prod.yml
services:
  backend:
    build: 
      context: ./backend      # Build context directory
      dockerfile: Dockerfile  # .dockerignore should be at ./backend/.dockerignore
```

### 2. Go Build Paths in Docker

**Working Configuration**:
- Build context: `./backend`
- WORKDIR: `/app`
- COPY: `. .` (copies backend/ contents to /app/)
- Go build: `./cmd/server/main.go` (explicit path to main.go file)

**Why explicit main.go path works better**:
- `./cmd/server` points to a directory, which can be ambiguous
- `./cmd/server/main.go` explicitly points to the buildable file
- More reliable across different Go versions and environments

## Deployment Gotchas

### 3. VPS Path Configuration

**Gotcha**: The deployment path must match between GitHub Actions workflow and VPS setup.

**Correct Configuration**:
- VPS path: `~/apps/chinese-app` (not `~/apps/chinese-graph`)
- GitHub Actions APP_DIR: `~/apps/chinese-app`
- Both must be identical for automated deployment to work

### 4. Docker Compose Version Compatibility

**Problem**: GitHub Actions failed with `docker-compose: command not found` on some VPS configurations.

**Solution**: Use fallback commands for both v1 and v2 syntax:
```bash
docker-compose -f docker-compose.prod.yml up -d --build || docker compose -f docker-compose.prod.yml up -d --build
```

**Better Solution**: Use only `docker compose` (v2) for consistency:
```bash
docker compose -f docker-compose.prod.yml up -d --build
```

### 5. Go Module Dependencies

**Problem**: Missing `go.sum` entries for updated dependencies causing build failures.

**Solution**: Always run `go mod tidy` after updating dependencies and commit the updated `go.sum`.

**Prevention**: Add this to your development workflow:
```bash
go mod download
go mod tidy
git add go.mod go.sum
```

## Environment Configuration

### 6. Environment Variables Structure

**Simplified approach**:
- Single `.env.example` file at repository root
- GitHub Actions creates `.env` automatically on VPS
- Docker Compose handles internal service communication via environment variables

**Production .env should contain**:
```bash
NODE_ENV=production
PORT=8080
NEO4J_PASSWORD=secure-password-here
```

## Architecture Decisions

### 7. Single Docker Compose Strategy

**Lesson**: Having multiple docker-compose files (development vs production) creates confusion.

**Solution**: Use single `docker-compose.prod.yml` for production with:
- Proper networking isolation
- Health checks for all services
- Persistent volumes for data
- Security best practices (non-root users)

### 8. Build Context Simplicity

**Best Practice**: Keep Dockerfile in the same directory as the code it builds.

**Structure**:
```
backend/
├── .dockerignore    # Docker ignore rules
├── Dockerfile       # Build instructions  
├── cmd/             # Application entry point
├── go.mod          # Go module definition
└── ...             # Source code
```

## GitHub Actions Best Practices

### 9. SSH Configuration

**Security Best Practice**: Use GitHub secrets for all sensitive data:
- `HOST`: VPS IP address
- `USER`: deployment user (e.g., 'deploy')
- `SSH_PRIVATE_KEY`: private key content
- `KNOWN_HOSTS`: output of `ssh-keyscan <host>`

### 10. Deployment Simplification

**Lesson**: Complex deployment scripts are error-prone.

**Simple deployment flow**:
1. SSH to VPS
2. Update git repository
3. Run `docker compose up -d --build`
4. Verify deployment

**Avoid**:
- Complex health checks in CI
- Multiple fallback strategies
- Overly verbose logging

## Go Specific Gotchas

### 11. Go Version Compatibility

**Problem**: Using non-existent Go versions (e.g., Go 1.25.0 when only 1.25.8 was available).

**Solution**: Always verify Docker image availability:
```bash
curl -s "https://registry.hub.docker.com/v2/repositories/library/golang/tags/"
```

**Use specific tags**: `golang:1.25.8-alpine3.23` instead of `golang:1.25-alpine`

### 12. Module-Aware Building

**Best Practice**: Always build with explicit paths in Docker:
```dockerfile
RUN go build -o server ./cmd/server/main.go
```

**Avoid**: Building directories without explicit main files:
```dockerfile
RUN go build -o server ./cmd/server  # Can be ambiguous
```

## Debugging Tips

### 13. Docker Build Debugging

**Add temporary debug steps** to Dockerfile when troubleshooting:
```dockerfile
COPY . .
RUN ls -la          # Show what was copied
RUN ls -la cmd/     # Check specific directories
```

**Remember to remove** debug steps after fixing the issue.

**Docker Cache Issues**: When troubleshooting build context issues, Docker layer caching can mask the problem. Even after fixing .dockerignore, cached layers may still use old build context. Add debug commands to invalidate cache and force rebuild.

**Context Transfer Size Indicator**: Watch for suspiciously small context transfer sizes (e.g., "transferring context: 794B" vs expected ~50KB+ for full backend). This indicates files are being excluded by .dockerignore patterns.

### 14. GitHub Actions Debugging

**Use detailed error logs**:
```bash
gh run view <run-id> --log-failed --repo <owner>/<repo>
```

**Check specific job output**:
```bash
gh run view <run-id> --job <job-id> --repo <owner>/<repo>
```

## Security Notes

### 15. Container Security

**Implemented best practices**:
- Non-root user in containers (`USER app`)
- No unnecessary port exposure
- Internal Docker network isolation
- Health checks for service monitoring

### 16. Secrets Management

**Never commit**:
- Real passwords or API keys
- Private keys or certificates
- Production environment files

**Always use**:
- GitHub Secrets for CI/CD
- Environment variables for runtime configuration
- `.env.example` files for documentation

## Performance Optimization

### 17. Docker Layer Caching

**Optimize Dockerfile layer order**:
```dockerfile
# 1. Install system dependencies (rarely changes)
RUN apk add --no-cache git

# 2. Copy dependency files (changes less frequently)
COPY go.mod go.sum ./

# 3. Download dependencies (cached when deps unchanged)
RUN go mod download

# 4. Copy source code (changes most frequently)
COPY . .

# 5. Build application
RUN go build ...
```

**This order maximizes Docker layer cache hits.**

## Maintenance Notes

### 18. Regular Updates

**Keep updated**:
- Go version in Dockerfile
- Dependencies in go.mod
- Base images (alpine, etc.)
- GitHub Actions versions

### 19. Monitoring

**Essential monitoring**:
- Application health endpoints (`/health`)
- Docker container status
- Disk space on VPS
- Service logs for errors

## Frontend & CI/CD Gotchas

### 20. NPM Package Lock Synchronization

**Problem**: `npm ci` fails with "package.json and package-lock.json are in sync" error after adding new dependencies.

**Error Message**:
```
npm error Missing: @testing-library/jest-dom@6.1.5 from lock file
npm error Missing: @testing-library/react@14.1.2 from lock file
```

**Root Cause**: Adding dependencies to `package.json` manually without updating `package-lock.json`.

**Solution**: 
1. Run `npm install` locally to regenerate lock file
2. Commit both `package.json` and `package-lock.json`
3. Keep using `npm ci` in CI for performance and reliability

**Prevention**: Always use `npm install <package>` instead of manually editing `package.json`.

### 21. ESLint Configuration with Next.js

**Problem**: ESLint fails with "Definition for rule '@typescript-eslint/*' was not found" when extending TypeScript ESLint configs.

**Root Cause**: `next/core-web-vitals` already includes TypeScript ESLint support, so manually adding `@typescript-eslint/recommended` creates conflicts.

**Solution**: Use only `next/core-web-vitals` and standard ESLint rules:
```json
{
  "extends": ["next/core-web-vitals"],
  "rules": {
    "prefer-const": "error"
  }
}
```

**Avoid**: Redundant TypeScript ESLint extensions when using Next.js presets.

### 22. GitHub Projects Classic Deprecation

**Problem**: GitHub Actions fail with "Projects (classic) is being deprecated" error.

**Failed Action**: `alex-page/github-project-automation-plus@v0.8.3`

**Solution**: Remove deprecated project automation workflows that use classic projects API.

**Alternative**: Use the new GitHub Projects experience with different automation or manual project management.

**Key Learning**: Monitor GitHub deprecation notices and update automation accordingly.

### 23. Testing Library Dependencies

**Problem**: TypeScript compilation fails with missing React Testing Library types.

**Missing Dependencies**:
- `@testing-library/react` - Component testing utilities
- `@testing-library/jest-dom` - Additional DOM matchers (`toBeInTheDocument`)

**Solution**: Add both dependencies to `devDependencies` and ensure Jest setup imports `@testing-library/jest-dom`.

**Configuration**: Jest setup should include:
```javascript
import '@testing-library/jest-dom'
```

---

## Quick Reference Commands

### Local Development
```bash
# Test Docker build locally
docker build -t chinese-graph-backend ./backend

# Run services locally  
docker compose -f docker-compose.prod.yml up -d

# Check logs
docker compose -f docker-compose.prod.yml logs -f backend
```

### Deployment Debugging
```bash
# Check latest deployment
gh run list --repo testprogmath/chinese-graph --limit 1

# View failed logs
gh run view <run-id> --log-failed --repo testprogmath/chinese-graph

# SSH to VPS for manual debugging
ssh deploy@204.168.217.208
cd ~/apps/chinese-app
docker compose -f docker-compose.prod.yml ps
```

### Go Module Management
```bash
# Update dependencies
go mod tidy

# Add new dependency
go get github.com/example/package

# Verify module integrity
go mod verify
```

---

*This document should be updated as new gotchas and lessons are discovered during development.*