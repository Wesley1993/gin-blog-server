-- 常用网站（友情链接）表
-- 幂等：CREATE TABLE IF NOT EXISTS + WHERE NOT EXISTS 种子，可重复执行
-- 适用于 PostgreSQL 16

CREATE TABLE IF NOT EXISTS blog_link (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    description VARCHAR(200) DEFAULT '',
    sort INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_link_status ON blog_link(status);
CREATE INDEX IF NOT EXISTS idx_link_sort ON blog_link(sort);

-- 种子数据：示例常用网站（按 name 判重）
INSERT INTO blog_link (name, url, description, sort, status, create_time, update_time)
SELECT 'Go 官网', 'https://go.dev', 'Go 语言官方网站', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_link WHERE name = 'Go 官网');

INSERT INTO blog_link (name, url, description, sort, status, create_time, update_time)
SELECT 'Gin 文档', 'https://gin-gonic.com', 'Gin Web 框架官方文档', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_link WHERE name = 'Gin 文档');

INSERT INTO blog_link (name, url, description, sort, status, create_time, update_time)
SELECT 'GitHub', 'https://github.com', '全球最大的代码托管平台', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_link WHERE name = 'GitHub');

INSERT INTO blog_link (name, url, description, sort, status, create_time, update_time)
SELECT 'V2EX', 'https://www.v2ex.com', '创意工作者们的社区', 4, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_link WHERE name = 'V2EX');
