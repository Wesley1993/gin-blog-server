#!/usr/bin/env bash
#
# Z-blog 数据备份脚本
# 用法：./scripts/backup.sh（在服务器上通过 cron 执行）
# cron: 15 3 * * * /opt/zblog/scripts/backup.sh >> /opt/zblog/backup/backup.log 2>&1
#
set -euo pipefail

BACKUP_DIR="/opt/zblog/backup"
KEEP_DAYS=7
STAMP=$(date +%F_%H%M)

cd /opt/zblog

echo "[$(date)] Starting backup..."

# 1. PostgreSQL logical backup
echo "[$(date)] Backing up PostgreSQL..."
docker compose exec -T postgres \
    pg_dump -U "${POSTGRES_USER:-appuser}" -d "${POSTGRES_DB:-blog}" -Fc \
    > "${BACKUP_DIR}/blog_${STAMP}.dump"

# 2. Configuration archive
echo "[$(date)] Archiving configuration..."
tar czf "${BACKUP_DIR}/conf_${STAMP}.tgz" \
    config.yaml .env nginx/conf.d docker-compose.yml 2>/dev/null || true

# 3. Cleanup old backups
echo "[$(date)] Cleaning backups older than ${KEEP_DAYS} days..."
find "$BACKUP_DIR" -name '*.dump' -mtime +$KEEP_DAYS -delete
find "$BACKUP_DIR" -name '*.tgz' -mtime +$KEEP_DAYS -delete

# 4. Report
LATEST_DUMP=$(ls -t "${BACKUP_DIR}"/blog_*.dump 2>/dev/null | head -1)
if [ -n "$LATEST_DUMP" ]; then
    SIZE=$(du -h "$LATEST_DUMP" | cut -f1)
    echo "[$(date)] Backup complete: $LATEST_DUMP ($SIZE)"
else
    echo "[$(date)] ERROR: No dump file created!"
    exit 1
fi
