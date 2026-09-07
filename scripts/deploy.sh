#!/usr/bin/env bash
#
# Z-blog 混合部署脚本（二进制 + Docker 基础设施）
# 用途：Jenkins Pipeline 和手工部署共用
# 用法：./scripts/deploy.sh [DEPLOY_HOST] [DEPLOY_USER]
#
set -euo pipefail

# ============ 配置 ============
DEPLOY_HOST="${1:-${DEPLOY_HOST:-}}"
DEPLOY_USER="${2:-${DEPLOY_USER:-root}}"
REMOTE_DIR="/opt/zblog"
TAG="$(date +%Y%m%d-%H%M%S)"
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
    error "DEPLOY_HOST is required."
    echo "Usage: $0 [DEPLOY_HOST] [DEPLOY_USER]"
    echo "  Or: export DEPLOY_HOST=zblog && $0"
    exit 1
fi

SSH_CMD="ssh -o StrictHostKeyChecking=accept-new -o ServerAliveInterval=15 ${DEPLOY_USER}@${DEPLOY_HOST}"
RSYNC_SSH="ssh -o StrictHostKeyChecking=accept-new -o ServerAliveInterval=15"

info "========================================="
info "  Z-blog 混合部署"
info "========================================="
info "  TAG:         $TAG"
info "  DEPLOY_HOST: $DEPLOY_HOST"
info "  DEPLOY_USER: $DEPLOY_USER"
info "  GOARCH:      $GOARCH"
info "  模式:        二进制 + Docker基础设施"
echo ""

# ============ Step 1: 确保远程目录存在 ============
info "Checking remote directory structure..."
$SSH_CMD "mkdir -p ${REMOTE_DIR}/{nginx/{conf.d,certs,logs},www/{blog,admin},backup,scripts}"
info "Remote directories ready"

# ============ Step 2: Build Backend ============
info "Building Go backend (linux/$GOARCH)..."
CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build -trimpath -ldflags="-s -w" -o bin/blog-server ./cmd
info "Backend built: bin/blog-server ($(du -h bin/blog-server | cut -f1))"

# ============ Step 3: Build Frontends ============
info "Building blog frontend..."
(cd blog && npm ci --silent && npm run build)
info "Blog frontend built"

info "Building admin frontend..."
(cd web && npm ci --silent && npm run build)
info "Admin frontend built"

# ============ Step 4: Sync infrastructure files ============
info "Syncing infrastructure files..."
rsync -az -e "$RSYNC_SSH" docker-compose.infra.yml "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/docker-compose.infra.yml"
rsync -az -e "$RSYNC_SSH" nginx/conf.d/ "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/nginx/conf.d/"
rsync -az -e "$RSYNC_SSH" scripts/backup.sh "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/scripts/backup.sh"
info "Infrastructure files synced"

# ============ Step 5: Ship binary + migrations + frontend ============
info "Shipping application files..."

# Go 二进制
rsync -az --progress -e "$RSYNC_SSH" bin/blog-server "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/blog-server.new"

# Migrations
rsync -az --delete -e "$RSYNC_SSH" migrations/ "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/migrations/"

# Frontend assets
rsync -az --delete -e "$RSYNC_SSH" blog/dist/ "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/www/blog/"
rsync -az --delete -e "$RSYNC_SSH" web/dist/ "${DEPLOY_USER}@${DEPLOY_HOST}:${REMOTE_DIR}/www/admin/"

info "Application files shipped"

# ============ Step 6: Activate on server ============
info "Activating new version on server..."
$SSH_CMD << EOF
set -euo pipefail
cd ${REMOTE_DIR}

echo "[远程] 备份数据库..."
pg_dump -h 127.0.0.1 -U \${POSTGRES_USER:-appuser} -Fc \${POSTGRES_DB:-blog} \
    > backup/pre_deploy_${TAG}.dump 2>/dev/null || echo "[远程] 警告: 数据库备份失败（可能首次部署）"

echo "[远程] 切换二进制..."
if [ -f blog-server ]; then
    cp blog-server blog-server.old
fi
mv blog-server.new blog-server
chmod +x blog-server

echo "[远程] 重启应用服务..."
sudo systemctl restart zblog

echo "[远程] 重载 Nginx..."
docker compose -f docker-compose.infra.yml exec nginx nginx -t 2>/dev/null && \
docker compose -f docker-compose.infra.yml exec nginx nginx -s reload || true

echo "[远程] 健康检查..."
sleep 5
if curl -fsS http://127.0.0.1:3000/health > /dev/null 2>&1; then
    echo "[远程] ✅ 健康检查通过"
    rm -f blog-server.old
else
    echo "[远程] ❌ 健康检查失败，回滚..."
    if [ -f blog-server.old ]; then
        mv blog-server.old blog-server
        sudo systemctl restart zblog
        sleep 3
    fi
    exit 1
fi
EOF

# ============ Step 7: Cleanup ============
info "Cleaning up local build artifacts..."
rm -rf bin/blog-server

echo ""
info "========================================="
info "  ✅ 部署成功！"
info "  TAG: $TAG"
info "  Host: $DEPLOY_HOST"
info "========================================="
