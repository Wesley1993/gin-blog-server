#!/usr/bin/env bash
#
# Z-blog 部署脚本
# 用途：Jenkins Pipeline 和手工部署共用
# 用法：./scripts/deploy.sh [TAG] [DEPLOY_HOST] [DEPLOY_USER]
#
set -euo pipefail

# ============ 配置 ============
TAG="${1:-$(date +%Y%m%d)-$(git rev-parse --short HEAD 2>/dev/null || echo 'manual')}"
DEPLOY_HOST="${2:-${DEPLOY_HOST:-}}"
DEPLOY_USER="${3:-${DEPLOY_USER:-deploy}}"
REMOTE_DIR="/opt/zblog"
GOARCH="${GOARCH:-amd64}"
GOOS="${GOOS:-linux}"

# ============ 颜色输出 ============
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; }

# ============ 前置检查 ============
if [ -z "$DEPLOY_HOST" ]; then
    error "DEPLOY_HOST is required. Usage: $0 [TAG] [DEPLOY_HOST] [DEPLOY_USER]"
    exit 1
fi

SSH_CMD="ssh -o StrictHostKeyChecking=accept-new -o ServerAliveInterval=15 ${DEPLOY_USER}@${DEPLOY_HOST}"

info "Deploying Z-blog"
info "  TAG:         $TAG"
info "  DEPLOY_HOST: $DEPLOY_HOST"
info "  DEPLOY_USER: $DEPLOY_USER"
info "  GOARCH:      $GOARCH"
echo ""

# ============ Step 1: Build Backend ============
info "Building Go backend (linux/$GOARCH)..."
CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build -trimpath -ldflags="-s -w" -o bin/blog-server ./cmd
info "Backend built: bin/blog-server ($(du -h bin/blog-server | cut -f1))"

# ============ Step 2: Build Frontends ============
info "Building blog frontend..."
(cd blog && npm ci --silent && npm run build)
info "Blog frontend built: blog/dist/"

info "Building admin frontend..."
(cd web && npm ci --silent && npm run build)
info "Admin frontend built: web/dist/"

# ============ Step 3: Build Docker Image ============
info "Building Docker image zblog-api:$TAG..."
docker build --platform linux/$GOARCH -f deploy/Dockerfile.api -t "zblog-api:$TAG" .
info "Docker image built"

# ============ Step 4: Ship to Server ============
info "Shipping image to $DEPLOY_HOST..."
docker save "zblog-api:$TAG" | gzip | $SSH_CMD "gunzip | docker load"
info "Image shipped"

info "Syncing frontend assets..."
rsync -az --delete blog/dist/ "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/www/blog/"
rsync -az --delete web/dist/  "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/www/admin/"
info "Frontend assets synced"

# ============ Step 5: Activate ============
info "Activating new version..."
$SSH_CMD << EOF
set -euo pipefail
cd $REMOTE_DIR

# Pre-deploy backup
echo "Backing up database..."
docker compose exec -T postgres pg_dump -U \${POSTGRES_USER:-appuser} -Fc \${POSTGRES_DB:-blog} > backup/pre_deploy_\$(date +%F_%H%M).dump 2>/dev/null || true

# Switch tag
sed -i.bak "s/^TAG=.*/TAG=$TAG/" .env

# Restart API (triggers migration)
docker compose up -d api

# Reload Nginx
docker compose exec nginx nginx -t 2>/dev/null
docker compose exec nginx nginx -s reload

# Health check (dual verification)
echo "Waiting for health check..."
sleep 10
# Verify API container is healthy
docker compose exec -T api wget -qO- http://127.0.0.1:3000/health > /dev/null 2>&1 || {
    echo "API container health check FAILED - rolling back"
    cp .env.bak .env
    docker compose up -d api
    sleep 5
    exit 1
}
# Verify Nginx proxy is working
curl -fsS http://127.0.0.1/health > /dev/null 2>&1 || {
    echo "Nginx proxy health check FAILED - rolling back"
    cp .env.bak .env
    docker compose up -d api
    sleep 5
    exit 1
}
echo "Health check PASSED"
EOF

# ============ Step 6: Cleanup ============
info "Cleaning up local Docker image..."
docker rmi "zblog-api:$TAG" 2>/dev/null || true
rm -rf bin/blog-server

echo ""
info "========================================="
info "  Deployment successful!"
info "  TAG: $TAG"
info "  Host: $DEPLOY_HOST"
info "========================================="
