package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"wuzhispace.com/pkg/logger"
)

// Run 按文件名顺序执行 migrations 目录下的所有 .sql 文件
// 使用 PostgreSQL 的 schema_migrations 表记录已执行的迁移，避免重复执行。
// 每个迁移文件在单个事务中执行，成功后写入版本记录。
func Run(db *sql.DB, dir string) error {
	// 1. 创建 schema_migrations 表（如不存在）
	const createTableSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("创建 schema_migrations 表失败: %w", err)
	}

	// 2. 读取目录下所有 .sql 文件，按文件名排序
	files, err := listSQLFiles(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		logger.Info("数据库迁移: 未找到迁移文件", "dir", dir)
		return nil
	}

	// 3. 逐个执行未应用过的迁移
	applied := 0
	for _, file := range files {
		version := filepath.Base(file)

		var exists int
		err := db.QueryRow("SELECT COUNT(1) FROM schema_migrations WHERE version = $1", version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("查询迁移记录失败 %s: %w", version, err)
		}
		if exists > 0 {
			logger.Debugf("数据库迁移: 跳过已执行的迁移 %s", version)
			continue
		}

		if err := applyMigration(db, file, version); err != nil {
			return fmt.Errorf("执行迁移 %s 失败: %w", version, err)
		}
		applied++
		logger.Info("数据库迁移: 执行成功", "version", version)
	}

	logger.Infof("数据库迁移完成: 本次执行 %d 个，共 %d 个迁移文件", applied, len(files))
	return nil
}

// listSQLFiles 列出目录下所有 .sql 文件，并按文件名排序
func listSQLFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取迁移目录 %s 失败: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".sql") {
			continue
		}
		files = append(files, filepath.Join(dir, name))
	}
	sort.Strings(files)
	return files, nil
}

// applyMigration 在单个事务中执行迁移文件内容，并记录版本号
func applyMigration(db *sql.DB, file, version string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("读取迁移文件失败: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(string(content)); err != nil {
		return fmt.Errorf("执行 SQL 失败: %w", err)
	}

	if _, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("记录迁移版本失败: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}
