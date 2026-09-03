# Blog 后台管理系统 — 需求清单（AI 编码参考手册）

> 基于 PRD v0.0.1 与现有代码脚手架分析生成，供 AI 分阶段编码使用。

---

## 1. 项目概述

| 维度 | 说明 |
|------|------|
| **项目名称** | Blog 后台管理系统 |
| **目标** | 博客后台管理：RBAC 权限体系、用户/角色/菜单管理、文章 CRUD + ES 全文检索、分类管理、站点配置、OSS 文件上传 |
| **技术栈** | Gin + GORM + PostgreSQL 16 + Redis + Elasticsearch 8.x + 阿里云 OSS SDK |
| **前端** | Vue3 + TS + Element-Plus（本项目仅后端） |
| **非目标** | ① 不开发前台 ② 不做第三方登录 ③ 文件不落地本地 ④ 不实现评论/访问统计 ⑤ ES 仅文章检索，PG 为真值源 |

---

## 2. 当前项目状态分析

### 2.1 已有什么（脚手架清单）

| 模块 | 文件 | 现状 |
|------|------|------|
| 入口 | `cmd/main.go` | 仅 `fmt.Println`，未接入 Gin |
| 配置 | `config/config.go` | 仅初始化 viper，未定义配置结构体 |
| 配置 | `config/config.yaml` | 含大量无关配置（rabbitmq、order、mall 前缀），需清理 |
| 路由 | `router/router.go` | 空路由，仅 Logger + Recovery 中间件 |
| 模型 | `internal/model/base.go` | BaseModel 使用 `gorm.DeletedAt`，需改为 `is_deleted smallint` 手动管理 |
| 模型 | `internal/model/article.go` | 字段与 PRD 不匹配（有 Author/Category/IsPublished，缺 CategoryID/Cover/Status/IsDeleted） |
| DTO | `internal/dto/article.go` | 仅有 CreateArticleReq，字段与 PRD 不匹配 |
| Handler | `internal/handler/article.go` | 空结构体，无方法 |
| Service | `internal/service/article.go` | 仅 CreateArticle，逻辑简单 |
| Service | `internal/service/base.go` | BaseService 含 DB + Redis |
| Repository | `internal/repository/article_repo.go` | 基础 CRUD，Delete 方法有 bug（`WHERE user_id` 无条件值） |
| JWT | `pkg/jwt/jwt.go` | 已实现 GenerateToken / ParseToken，Claims 含多余字段（ClientType） |
| Redis | `pkg/redis/redis.go` | 封装完整，支持常用操作 |
| Response | `pkg/response/response.go` | 统一响应 + 泛型分页，**Success 返回 code=0，PRD 要求 code=200** |
| 错误码 | `pkg/errors/code.go` | 已有基础错误码体系，但错误码值与 PRD 定义不一致 |
| 验证码 | `pkg/captcha/captcha.go` | 滑动验证码，已集成 |
| 工具 | `pkg/util/util.go` | PageQuery + 泛型 Paginate，可用 |
| 依赖 | `go.mod` | Gin、GORM、go-redis/v8、viper、jwt/v5、captcha 已引入；缺 PG 驱动、bcrypt、ES、OSS SDK |

### 2.2 缺什么（需新建的模块）

| 模块 | 说明 |
|------|------|
| **数据模型** | `sys_user`、`sys_role`、`sys_menu`、`blog_category`、`blog_site_config` 模型及 DTO |
| **认证模块** | login / logout / userinfo 的 handler → service → repository |
| **RBAC 中间件** | JWT 鉴权中间件 + RBAC 权限校验中间件 |
| **限流中间件** | 基于 Redis 的登录接口限流（5次/分钟/IP） |
| **用户管理** | handler / service / repository 全套 |
| **角色管理** | handler / service / repository 全套 |
| **菜单管理** | handler / service / repository 全套 |
| **分类管理** | handler / service / repository 全套 |
| **站点配置** | handler / service / repository + Redis 缓存 |
| **ES 客户端** | `pkg/es/` — go-elasticsearch/v8 封装、索引管理、同步、搜索 |
| **OSS 客户端** | `pkg/oss/` — 阿里云 OSS SDK 封装、上传、连通测试 |
| **数据库初始化** | `cmd/main.go` 中 PG/Redis/ES 初始化流程 |
| **建表 SQL** | `migrations/` 目录下 6 张表 DDL + 索引 + 初始数据 |
| **GORM PG 驱动** | `gorm.io/driver/postgres` 尚未引入 |
| **bcrypt 依赖** | `golang.org/x/crypto/bcrypt` 用于密码加密 |

### 2.3 需要修正什么（与 PRD 不匹配）

| 问题 | 当前 | 应改为 |
|------|------|--------|
| BaseModel 软删除 | `gorm.DeletedAt` 自动软删 | 移除 DeletedAt，改用 `IsDeleted int16` 手动管理 |
| BaseModel 字段 | 含 Remarks 字段；ID 为 `uint` | 移除 Remarks；ID 改为 `int64`（对应 bigserial） |
| BaseModel 时间字段 | `CreatedAt`/`UpdatedAt` | 改为 `CreateTime`/`UpdateTime`，对应 PG `create_time`/`update_time` |
| Article 模型 | 有 Author/Category/IsPublished/IsDeleted(bool) | 改为 CategoryID(int64)/Cover/Status(int16)/IsDeleted(int16) |
| config.yaml | 含 rabbitmq/order/upload/mall 前缀等无关配置 | 清理为 blog 专用配置；OSS 参数不入 yaml；JWT prefix 改为 `blog_` |
| config.go | 无结构体定义 | 定义完整 Config struct 并绑定 viper |
| 错误码 | 错误码值体系（文章10000+、用户20000+…） | 按 PRD 重新定义：10001-10005 用户，20001-20002 分类，30001 文章，40001-40004 OSS，50001 ES |
| JWT Claims | 含 ClientType 等多余字段 | 精简为 UserID/Username/RoleID/IsSuper |
| response.go | Success 返回 code=0 | PRD 定义成功 code=200，需统一 |
| article_repo Delete | `WHERE user_id` 缺少条件值，有 bug | 文章为逻辑删除，改为 `UPDATE is_deleted=1` |

---

## 3. 关键技术决策

| 决策项 | 选型 | 备注 |
|--------|------|------|
| ES 客户端 | `github.com/elastic/go-elasticsearch/v8` | 官方客户端 |
| OSS SDK | `github.com/aliyun/alibabacloud-oss-go-sdk-v2` | 阿里云官方 v2 SDK，V4 签名 |
| Redis 客户端 | `github.com/go-redis/redis/v8` | 保持现有 |
| PG 驱动 | `gorm.io/driver/postgres` | 需新增依赖 |
| 密码加密 | `golang.org/x/crypto/bcrypt` | 需新增依赖 |
| OSS 配置来源 | 运行时从 PG `blog_site_config` 表动态读取 | 不写死在 config.yaml |
| ES 降级 | ES 不可用时自动降级为 PG `LIKE` 模糊查询 | 日志打印 ES 异常 |
| BaseModel | 移除 `gorm.DeletedAt`，改用 `is_deleted smallint` 手动管理 | 所有查询需手动加 `WHERE is_deleted = 0` |
| 主键类型 | `bigserial`（Go 侧 `int64`） | PRD 约定 |
| 时间字段 | `timestamp without time zone`（Go 侧 `time.Time`，JSON tag `create_time`/`update_time`） | 不使用 gorm 自动时间戳 |
| JSON 字段 | `jsonb`（如 `menu_ids`、`button_perms`） | 使用 `datatypes.JSON` 或自定义序列化 |
| ES 分词 | 索引用 `ik_max_word`，搜索用 `ik_smart` | mapping 中同时指定 analyzer 和 search_analyzer |
| OSS 文件命名 | `blog/{year}/{month}/{uuid}.{ext}` | 按日期分目录，UUID 防冲突 |
| OSS 上传限制 | 单文件最大 10MB；仅允许 `image/jpeg`、`image/png`、`image/gif`、`image/webp` | 超出返回 40003/40004 错误码 |
| OSS 连通测试 | 上传 4 字节测试文件后立即删除 | 同时验证读写权限 |
| OSS 接口抽象 | 定义 `ObjectStorage` 接口（Upload/Delete/TestConnectivity） | 当前实现 AliyunOSS，未来可扩展 S3 |
| JWT token prefix | `blog_` | PRD 约定 |

---

## 4. 约束与规范（AI 编码必须遵守）

> 以下为全局编码约束，所有 Phase 的实现均须遵守。

### 4.1 时间处理

- **禁止**在 GORM model 中使用 `CreatedAt`/`UpdatedAt` 自动时间戳字段
- **禁止**使用 `gorm.DeletedAt` 自动软删除
- 时间字段统一使用 `time.Time` 类型，JSON tag 为 `create_time` / `update_time`
- 时间的创建和更新在 service 层手动赋值（`time.Now()`），不依赖 GORM 自动管理
- PG 侧类型为 `timestamp without time zone`

### 4.2 逻辑删除

- **禁止**使用 `gorm.DeletedAt` 或 GORM 的自动软删除机制
- 所有需要逻辑删除的表使用 `is_deleted int16` 字段（`smallint NOT NULL DEFAULT 0`，0=未删 1=已删）
- 所有查询必须手动添加 `WHERE is_deleted = 0` 条件
- 删除操作统一为 `UPDATE ... SET is_deleted = 1`

### 4.3 主键类型

- 主键统一使用 `int64` 类型（对应 PG `bigserial`）
- **禁止**使用 `uint` 或 `uint64` 作为主键类型

### 4.4 依赖注入

- **禁止**使用 `init()` 函数进行初始化
- 所有依赖通过构造函数或显式函数参数注入
- 示例：`func NewArticleService(db *gorm.DB, redis *redis.RedisClient) *ArticleService`

### 4.5 配置管理

- 所有配置必须定义对应的 Go struct，通过 viper 绑定
- **禁止**在代码中硬编码配置值（如 Redis key 前缀、JWT secret 等）
- 配置值统一从 `config.Config` 结构体读取

### 4.6 错误处理

- 所有错误必须显式处理，**禁止**使用 `_` 忽略 error
- 错误返回使用 `pkg/errors` 中定义的错误码和消息
- 数据库操作、Redis 操作、ES 操作等均需检查并处理 error

### 4.7 代码风格

- 遵循 Go 标准代码风格
- 包名使用小写短单词（如 `model`、`service`、`handler`）
- 导出函数/类型必须有注释
- 接口命名使用 `-er` 后缀或 `I` 前缀（如 `ObjectStorage`）

---

## 5. 数据模型清单

### 5.0 公共基础模型

所有表共有字段：

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| create_time | timestamp without time zone | NOT NULL | 创建时间 |
| update_time | timestamp without time zone | NOT NULL | 更新时间 |

> 逻辑删除说明：所有表使用应用层管理 `is_deleted` 字段（`smallint NOT NULL DEFAULT 0`），**不使用 GORM 自动软删除**。所有查询需手动加 `WHERE is_deleted = 0`。

### 5.1 sys_user 用户表 🆕

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| username | varchar(50) | NOT NULL, UNIQUE | 登录账号 |
| password | varchar(100) | NOT NULL | bcrypt 加密密码 |
| nickname | varchar(50) | | 昵称 |
| role_id | bigint | NOT NULL | 绑定角色 ID |
| status | smallint | NOT NULL, DEFAULT 1 | 0 禁用 1 启用 |
| create_time | timestamp | NOT NULL | |
| update_time | timestamp | NOT NULL | |

> **差异**：全新表，现有代码无。

### 5.2 sys_role 角色表 🆕

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| role_name | varchar(50) | NOT NULL | 角色名称 |
| menu_ids | jsonb | | 菜单 ID 数组 |
| button_perms | jsonb | | 按钮权限标识数组，如 `["article:add","article:edit"]` |
| is_super | smallint | NOT NULL, DEFAULT 0 | 是否超级管理员 0/1 |
| create_time | timestamp | NOT NULL | |
| update_time | timestamp | NOT NULL | |

> **差异**：全新表。相比初版 PRD 新增 `button_perms` 字段，将按钮权限标识独立存储。

### 5.3 sys_menu 菜单权限表 🆕

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| parent_id | bigint | NOT NULL, DEFAULT 0 | 父菜单 ID |
| menu_name | varchar(50) | NOT NULL | 菜单名称 |
| menu_type | smallint | NOT NULL | 1 目录 2 页面 3 按钮 |
| path | varchar(100) | | 前端路由 path |
| perms | varchar(100) | | 权限标识如 `article:add` |
| sort | int | NOT NULL, DEFAULT 0 | 排序值 |
| status | smallint | NOT NULL, DEFAULT 1 | 0 禁用 1 启用 |
| create_time | timestamp | NOT NULL | |
| update_time | timestamp | NOT NULL | |

> **差异**：全新表。相比初版 PRD 新增 `status`、`create_time`、`update_time` 字段。

### 5.4 blog_category 分类表 🆕

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| name | varchar(60) | NOT NULL | 分类名称 |
| parent_id | bigint | NOT NULL, DEFAULT 0 | 父分类 ID |
| sort | int | NOT NULL, DEFAULT 0 | 排序 |
| status | smallint | NOT NULL, DEFAULT 1 | 0 禁用 1 启用 |
| create_time | timestamp | NOT NULL | |
| update_time | timestamp | NOT NULL | |

> **差异**：全新表。

### 5.5 blog_article 文章表 ⚠️ 需重构

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| title | varchar(120) | NOT NULL | 标题 |
| category_id | bigint | NOT NULL | 分类 ID |
| cover | varchar(255) | | 封面 OSS 地址 |
| content | text | NOT NULL | 富文本正文 |
| tags | varchar(200) | | 标签逗号分隔 |
| status | smallint | NOT NULL, DEFAULT 0 | 0 草稿 1 已发布 |
| is_deleted | smallint | NOT NULL, DEFAULT 0 | 0 未删 1 已删 |
| create_time | timestamp | NOT NULL | |
| update_time | timestamp | NOT NULL | |

**索引定义**：

```sql
CREATE INDEX idx_article_category ON blog_article(category_id);
CREATE INDEX idx_article_status ON blog_article(status);
CREATE INDEX idx_article_create_time ON blog_article(create_time);
```

> **与现有代码差异**：移除 Author、Category(string)、IsPublished(bool)、IsDeleted(bool)；新增 CategoryID(int64)、Cover、Status(int16)、IsDeleted(int16)。新增 3 个索引。

### 5.6 blog_site_config 站点配置表 🆕

| 字段 | PG 类型 | 约束 | 说明 |
|------|---------|------|------|
| id | bigserial | PK | 主键 |
| site_name | varchar(100) | | 站点名称 |
| site_desc | varchar(255) | | 站点描述 |
| copyright | varchar(255) | | 版权 |
| oss_access_key | varchar(200) | | OSS AK |
| oss_secret_key | varchar(200) | | OSS SK |
| oss_bucket | varchar(100) | | OSS Bucket |
| oss_endpoint | varchar(200) | | OSS Endpoint |
| oss_domain | varchar(200) | | OSS 自定义域名 |
| create_time | timestamp | NOT NULL | |
| update_time | timestamp | NOT NULL | |

> **差异**：全新表。OSS 配置存在此表中，运行时动态读取。

### 5.7 ES 索引：`blog_article_index` 🆕

```json
{
  "mappings": {
    "properties": {
      "id":          { "type": "long" },
      "title":       { "type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart" },
      "content":     { "type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart" },
      "tags":        { "type": "keyword" },
      "category_id": { "type": "long" },
      "status":      { "type": "integer" },
      "create_time": { "type": "date" }
    }
  }
}
```

> **注意**：title/content 同时指定 `analyzer`（索引时 ik_max_word 最大切词提高召回）和 `search_analyzer`（搜索时 ik_smart 智能分词提高精度）。

---

## 6. API 接口清单

### 统一规范

**统一响应格式**：

```json
{ "code": 200, "msg": "ok", "data": {} }
```

> code=200 成功；401 未登录；403 无权限；其他为业务错误码。

**分页参数规范**：

- `page`：页码，从 1 开始，默认 1
- `page_size`：每页条数，默认 10，最大 100
- 分页响应：`{"code":200,"msg":"ok","data":{"list":[...],"total":100,"page":1,"page_size":10}}`

### 5.1 认证接口

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | `/api/auth/login` | 账号密码登录 | 无需登录（白名单） |
| POST | `/api/auth/logout` | 登出，清除 Redis token | 需登录 |
| GET | `/api/auth/userinfo` | 获取当前用户信息+菜单树+权限标识集合 | 需登录 |

**POST `/api/auth/login`**

请求：
```json
{ "username": "admin", "password": "123456" }
```
响应：
```json
{ "code": 200, "msg": "ok", "data": { "token": "eyJhbGci..." } }
```

**GET `/api/auth/userinfo`**

响应：
```json
{
  "code": 200, "msg": "ok",
  "data": {
    "user": { "id": 1, "username": "admin", "nickname": "管理员", "role_id": 1 },
    "menus": [{ "id": 1, "parent_id": 0, "menu_name": "权限管理", "menu_type": 1, "path": "/system", "children": [...] }],
    "permissions": ["article:add", "article:edit", "article:delete"]
  }
}
```

### 5.2 RBAC 权限

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/menu/list` | 获取全部菜单树 | 需登录 |
| GET | `/api/role/list` | 角色列表 | 需登录 |
| POST | `/api/role/create` | 创建角色 | 需登录 |
| PUT | `/api/role/update` | 更新角色+权限 | 需登录 |
| DELETE | `/api/role/:id` | 删除角色（超管角色不可删） | 需登录 |

**POST `/api/role/create`**

请求：
```json
{ "role_name": "编辑", "menu_ids": [1, 2, 3], "button_perms": ["article:add", "article:edit"] }
```

**PUT `/api/role/update`**

请求：
```json
{ "id": 2, "role_name": "编辑", "menu_ids": [1, 2, 3], "button_perms": ["article:add", "article:edit"] }
```

**GET `/api/role/list`**

响应：
```json
{ "code": 200, "data": [{ "id": 1, "role_name": "超级管理员", "menu_ids": [1,2,3], "button_perms": ["article:add"], "is_super": 1 }] }
```

### 5.3 用户管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/user/page` | 用户分页 | 需登录 |
| POST | `/api/user/create` | 创建用户 | 需登录 |
| PUT | `/api/user/update` | 编辑用户 | 需登录 |
| PUT | `/api/user/resetPwd/:id` | 重置密码 | 需登录 |
| DELETE | `/api/user/:id` | 删除用户（超管不可删） | 需登录 |

**GET `/api/user/page`**

请求参数（Query）：`page=1&page_size=10`

响应：
```json
{ "code": 200, "data": { "list": [{ "id": 1, "username": "admin", "nickname": "管理员", "role_id": 1, "status": 1 }], "total": 50, "page": 1, "page_size": 10 } }
```

**POST `/api/user/create`**

请求：
```json
{ "username": "editor", "nickname": "编辑", "password": "123456", "role_id": 2, "status": 1 }
```

**PUT `/api/user/update`**

请求：
```json
{ "id": 2, "nickname": "编辑A", "role_id": 2, "status": 1 }
```

**PUT `/api/user/resetPwd/:id`**

请求：
```json
{ "password": "newpassword" }
```

### 5.4 分类管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/category/tree` | 获取分类树 | 需登录 |
| POST | `/api/category/create` | 新增分类 | 需登录 |
| PUT | `/api/category/update` | 更新分类 | 需登录 |
| DELETE | `/api/category/:id` | 删除分类（被引用拦截） | 需登录 |

**GET `/api/category/tree`**

响应：
```json
{ "code": 200, "data": [{ "id": 1, "name": "技术", "parent_id": 0, "sort": 1, "status": 1, "children": [{ "id": 2, "name": "Go", "parent_id": 1, "sort": 1, "status": 1 }] }] }
```

**POST `/api/category/create`**

请求：
```json
{ "name": "Go", "parent_id": 1, "sort": 1, "status": 1 }
```

**PUT `/api/category/update`**

请求：
```json
{ "id": 2, "name": "Golang", "parent_id": 1, "sort": 1, "status": 1 }
```

### 5.5 文章管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/article/page` | PG 普通分页列表 | 需登录 |
| GET | `/api/article/es/search` | ES 全文搜索 | 需登录 |
| POST | `/api/article/create` | 创建文章，内部同步 ES | 需登录 |
| PUT | `/api/article/update` | 更新文章，内部同步 ES | 需登录 |
| DELETE | `/api/article/:id` | 逻辑删除，同步删除 ES 文档 | 需登录 |
| POST | `/api/article/es/rebuild` | 手动重建 ES 全量索引 | 需登录 + 超级管理员 |

**GET `/api/article/page`**

请求参数（Query）：`page=1&page_size=10&category_id=1&status=1`

响应：
```json
{ "code": 200, "data": { "list": [{ "id": 1, "title": "Hello Go", "category_id": 1, "cover": "https://oss.../cover.jpg", "tags": "go,web", "status": 1, "create_time": "2025-01-01 00:00:00" }], "total": 100, "page": 1, "page_size": 10 } }
```

**GET `/api/article/es/search`**

请求参数（Query）：`keyword=Go开发&page=1&page_size=10`

响应：
```json
{ "code": 200, "data": { "list": [{ "id": 1, "title": "Hello Go", "tags": "go,web", "category_id": 1, "status": 1, "create_time": "2025-01-01 00:00:00" }], "total": 5, "page": 1, "page_size": 10 } }
```

**POST `/api/article/create`**

请求：
```json
{ "title": "Hello Go", "category_id": 1, "cover": "https://oss.../cover.jpg", "content": "<p>正文内容</p>", "tags": "go,web", "status": 1 }
```

**PUT `/api/article/update`**

请求：
```json
{ "id": 1, "title": "Hello Go v2", "category_id": 1, "content": "<p>更新后正文</p>", "tags": "go", "status": 1 }
```

### 5.6 站点 & OSS

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/site/config` | 获取站点配置 | 需登录 |
| PUT | `/api/site/config` | 保存配置，刷新 Redis 缓存 | 需登录 |
| POST | `/api/site/testOss` | OSS 连通测试 | 需登录 |
| POST | `/api/upload/oss` | form-data 文件上传，返回 url | 需登录 |

**GET `/api/site/config`**

响应：
```json
{ "code": 200, "data": { "site_name": "我的博客", "site_desc": "一个技术博客", "copyright": "©2025", "oss_access_key": "***", "oss_secret_key": "***", "oss_bucket": "my-blog", "oss_endpoint": "oss-cn-hangzhou.aliyuncs.com", "oss_domain": "cdn.example.com" } }
```

**PUT `/api/site/config`**

请求：
```json
{ "site_name": "我的博客", "site_desc": "一个技术博客", "copyright": "©2025", "oss_access_key": "LTAI...", "oss_secret_key": "xxx", "oss_bucket": "my-blog", "oss_endpoint": "oss-cn-hangzhou.aliyuncs.com", "oss_domain": "cdn.example.com" }
```

**POST `/api/upload/oss`**

请求：`Content-Type: multipart/form-data`，字段 `file`

> 限制：单文件最大 **10MB**；仅支持 `image/jpeg`、`image/png`、`image/gif`、`image/webp`。超出返回错误码 40003（类型不支持）或 40004（大小超限）。

响应：
```json
{ "code": 200, "data": { "url": "https://cdn.example.com/blog/2025/01/a1b2c3d4.png" } }
```

---

## 7. 中间件清单

### 6.1 JWT 鉴权中间件

**路径**：`middleware/jwt.go`

**流程**：
1. 从 `Authorization: Bearer {token}` 提取 token
2. 调用 `jwt.ParseToken` 校验 JWT 有效性
3. 校验 Redis 中 `token:{userId}` 是否存在（防止 token 被 logout 后继续使用）
4. 将 `userID`、`roleID`、`isSuper` 注入 `gin.Context`
5. 失败返回 401

**白名单路由**：`/api/auth/login`

### 6.2 RBAC 权限中间件

**路径**：`middleware/rbac.go`

**流程**：
1. 从 Context 取出当前用户 `roleID` 和 `isSuper`
2. 若 `is_super=1`（超级管理员），直接放行
3. 从 Redis 读取 `rbac:perm:{userId}` 权限标识集合
4. 缓存未命中则查 PG：通过 `role_id` → `sys_role.button_perms` 获取权限标识集合
5. 结果写入 Redis 缓存
6. 获取当前路由所需的权限标识（通过路由 metadata 或自定义 header）
7. 比对，不匹配返回 403

### 6.3 限流中间件

**路径**：`middleware/ratelimit.go`

**策略**：
- 登录接口专用：基于 Redis `INCR + EXPIRE` 实现
- Key 模式：`rate:login:{ip}`
- 限制：5 次/分钟/IP
- 超限返回 429 Too Many Requests

---

## 8. Redis 缓存设计

| Key 模式 | 类型 | TTL | 写入时机 | 失效策略 |
|----------|------|-----|----------|----------|
| `token:{userId}` | string | 24h | 登录时写入 | logout 时删除；自然过期 |
| `rbac:perm:{userId}` | string (JSON) | 24h | 首次鉴权时写入 | 角色修改时批量删除该角色下所有用户的 key；用户信息变更时删除 |
| `rbac:menu:{userId}` | string (JSON) | 24h | userinfo 接口首次查询时写入 | 同上 |
| `site:config` | string (JSON) | 1h | 获取站点配置时写入 | 保存配置时删除 |
| `rate:login:{ip}` | string (int) | 60s | 每次登录请求 INCR | 自然过期 |

### 缓存失效触发点

| 触发场景 | 操作 |
|----------|------|
| **角色修改**（菜单/权限变更） | 查询该角色下所有用户 → 批量删除其 `rbac:perm:{userId}` 和 `rbac:menu:{userId}` |
| **用户信息变更**（编辑/禁用/删除） | 删除该用户 `token:{userId}`、`rbac:perm:{userId}`、`rbac:menu:{userId}` |
| **站点配置保存** | 删除 `site:config` key |

---

## 9. ES 集成设计

### 8.1 索引

- 索引名：`blog_article_index`
- Mapping 见 §5.7
- 分词器：索引 `ik_max_word`（最大切词提高召回），搜索 `ik_smart`（智能分词提高精度）

### 8.2 IK 分词器依赖

ES 必须预装 IK 插件。安装命令示例：

```bash
./bin/elasticsearch-plugin install https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-8.17.0.zip
```

> 版本号需与 ES 版本一致。

### 8.3 同步策略

| 触发场景 | ES 操作 | 说明 |
|----------|---------|------|
| 文章创建 | upsert 文档 | PG 写入成功后同步 |
| 文章更新 | upsert 文档 | PG 更新成功后同步 |
| 文章逻辑删除 | delete 文档 | PG 标记 is_deleted=1 后，从 ES 删除（ES 不存逻辑删除数据） |
| 手动重建索引 | delete index + create + bulk upsert | 全量从 PG 读取 is_deleted=0 的文章 |

### 8.4 搜索实现

```
GET /blog_article_index/_search
{
  "query": {
    "multi_match": {
      "query": "keyword",
      "type": "best_fields",
      "fields": ["title^3", "content", "tags^2"]
    }
  },
  "highlight": {
    "fields": { "title": {}, "content": {} },
    "pre_tags": ["<em>"], "post_tags": ["</em>"]
  },
  "sort": [{ "create_time": { "order": "desc" } }],
  "from": 0, "size": 10
}
```

> - `type: "best_fields"`：每个 token 取所有字段中的最高分，避免同一字段多次匹配导致分数膨胀，适合短文本搜索场景。
> - 搜索结果按 `create_time` 降序排列（最新文章优先）。
> - 搜索仅查 ES；命中后如需详情，从 PG 读取。

### 8.5 降级方案

ES 不可用时，全文搜索降级为 PG 查询：

```sql
WHERE is_deleted=0 AND (title LIKE '%keyword%' OR content LIKE '%keyword%' OR tags LIKE '%keyword%')
```

日志打印 ES 异常。

### 8.6 启动时检查

1. 尝试 `ping` ES，验证连通性
2. 调用 `GET /_analyze` 验证 IK 分词器可用性：`{"analyzer":"ik_max_word","text":"测试中文分词"}`
3. 检查 `blog_article_index` 是否存在（`Indices.Exists`），不存在则自动创建（`Indices.Create`）
4. 以上检查失败**不阻止应用启动**，仅日志打印 ERROR，搜索自动降级

---

## 10. 错误码表

### 9.1 通用错误码

| 错误码 | 含义 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未登录 / token 过期 |
| 403 | 无权限 |
| 429 | 请求过于频繁（限流） |
| 500 | 服务器内部错误 |

### 9.2 业务错误码（PRD 定义）

| 错误码 | 说明 |
|--------|------|
| 10001 | 用户名已存在 |
| 10002 | 角色不存在 |
| 10003 | 超级管理员不可删除 |
| 10004 | 账号已禁用 |
| 10005 | 用户名或密码错误 |
| 20001 | 分类不存在 |
| 20002 | 分类被文章引用，禁止删除 |
| 30001 | 文章不存在 |
| 40001 | OSS 配置未设置 |
| 40002 | OSS 上传失败 |
| 40003 | 文件类型不支持 |
| 40004 | 文件大小超限 |
| 50001 | ES 搜索降级 |

---

## 11. 分阶段实现计划

### Phase 0：基础设施

| 编号 | 需求描述 | 涉及文件 | 依赖 | 验收标准 |
|------|----------|----------|------|----------|
| 0.1 | 清理 config.yaml：移除 rabbitmq/order/upload 无关配置；修正 app.name 为 `blog-admin`；JWT prefix 改为 `blog_`；ES 配置补充 index_name；移除 OSS 配置段（由 DB 动态管理） | `config/config.yaml` | 无 | 配置文件仅含 app/database/redis/jwt/elasticsearch/log 段 |
| 0.2 | 重构 config.go：定义完整 Config struct（App/Database/Redis/JWT/ES/Log），实现 `InitConfig() *Config` | `config/config.go` | 0.1 | 能正确解析 yaml 到结构体 |
| 0.3 | 添加依赖：`gorm.io/driver/postgres`、`golang.org/x/crypto`、`github.com/elastic/go-elasticsearch/v8`、`github.com/aliyun/alibabacloud-oss-go-sdk-v2` | `go.mod` | 无 | `go mod tidy` 成功 |
| 0.4 | 创建建表 SQL：6 张表 DDL + 文章表索引 + 初始超级管理员 + 初始菜单数据 | `migrations/001_init.sql` | 无 | SQL 可在 PG16 执行成功 |
| 0.5 | 重构 BaseModel：移除 DeletedAt/Remarks，ID 改为 `int64`，时间字段改为 `CreateTime`/`UpdateTime`（tag `create_time`/`update_time`） | `internal/model/base.go` | 0.3 | 所有模型使用新 BaseModel |
| 0.6 | 初始化 PG 连接：在 main.go 中初始化 GORM + PG | `cmd/main.go` | 0.2, 0.3 | 应用启动可连接 PG |
| 0.7 | 初始化 Redis 连接 | `cmd/main.go` | 0.2, 0.3 | 应用启动可连接 Redis |
| 0.8 | 修正 response.go：Success 返回 code=200（与 PRD 统一） | `pkg/response/response.go` | 无 | 所有成功响应 code=200 |
| 0.9 | 重构错误码：按 PRD 错误码表重写 | `pkg/errors/code.go` | 无 | 错误码与 §9 一致 |

### Phase 1：认证与 RBAC

| 编号 | 需求描述 | 涉及文件 | 依赖 | 验收标准 |
|------|----------|----------|------|----------|
| 1.1 | 创建 `sys_user`、`sys_role`、`sys_menu` 模型 | `internal/model/user.go`、`role.go`、`menu.go` | 0.5 | 模型字段与 §5 一致 |
| 1.2 | 创建认证 DTO（LoginReq/LoginResp/UserInfoResp） | `internal/dto/auth.go` | 1.1 | |
| 1.3 | 实现 UserRepository（FindByUsername、FindByID、FindByRoleID） | `internal/repository/user_repo.go` | 1.1 | |
| 1.4 | 实现 RoleRepository（FindByID、FindAll、FindUsersByRoleID） | `internal/repository/role_repo.go` | 1.1 | |
| 1.5 | 实现 MenuRepository（FindAll、FindByIDs、BuildTree） | `internal/repository/menu_repo.go` | 1.1 | |
| 1.6 | 实现 AuthService（Login、Logout、GetUserInfo） | `internal/service/auth.go` | 1.2-1.5 | 登录返回 JWT token |
| 1.7 | 实现 AuthHandler（login/logout/userinfo） | `internal/handler/auth.go` | 1.6 | 接口可调通 |
| 1.8 | 实现 JWT 鉴权中间件 | `middleware/jwt.go` | 0.7 | 无 token 返回 401 |
| 1.9 | 实现 RBAC 权限中间件 | `middleware/rbac.go` | 1.8 | 无权限返回 403 |
| 1.10 | 实现限流中间件（5次/分钟/IP） | `middleware/ratelimit.go` | 0.7 | 超限返回 429 |
| 1.11 | 用户管理 CRUD（handler/service/repo/dto） | `internal/handler/user.go`、`internal/service/user.go`、`internal/dto/user.go` | 1.8, 1.9 | 分页/创建/编辑/重置密码/删除可用 |
| 1.12 | 角色管理 CRUD | `internal/handler/role.go`、`internal/service/role.go`、`internal/dto/role.go` | 1.8, 1.9 | 列表/创建/更新/删除可用，超管角色不可删 |
| 1.13 | 菜单管理（列表树） | `internal/handler/menu.go`、`internal/service/menu.go` | 1.8, 1.9 | 菜单树接口可用 |
| 1.14 | Redis 缓存层：token 存储、权限缓存、缓存失效逻辑 | `internal/service/cache.go` 或集成到各 service | 1.8 | 权限缓存命中/失效正常 |

### Phase 2：业务模块

| 编号 | 需求描述 | 涉及文件 | 依赖 | 验收标准 |
|------|----------|----------|------|----------|
| 2.1 | 重构 Article 模型，与 §5.5 对齐（含索引 tag） | `internal/model/article.go` | 0.5 | 字段完全匹配 |
| 2.2 | 创建 `blog_category` 模型 + DTO | `internal/model/category.go`、`internal/dto/category.go` | 0.5 | |
| 2.3 | 创建 `blog_site_config` 模型 + DTO | `internal/model/site_config.go`、`internal/dto/site.go` | 0.5 | |
| 2.4 | 重构 Article DTO | `internal/dto/article.go` | 2.1 | 与 API 定义一致 |
| 2.5 | 分类管理 CRUD（handler/service/repo） | `internal/handler/category.go`、`internal/service/category.go`、`internal/repository/category_repo.go` | 2.2 | 树/创建/更新/删除可用，被引用分类不可删 |
| 2.6 | 文章管理 CRUD（修正，不含 ES 同步） | `internal/handler/article.go`、`internal/service/article.go`、`internal/repository/article_repo.go` | 2.1, 2.4 | PG 分页/创建/更新/逻辑删除可用 |
| 2.7 | 站点配置管理（handler/service/repo） | `internal/handler/site.go`、`internal/service/site.go`、`internal/repository/site_config_repo.go` | 2.3 | 获取/保存配置可用 |
| 2.8 | 站点配置 Redis 缓存（保存时删除 `site:config`） | 集成到 site service | 2.7 | 保存后缓存刷新 |

### Phase 3：外部集成

| 编号 | 需求描述 | 涉及文件 | 依赖 | 验收标准 |
|------|----------|----------|------|----------|
| 3.1 | ES 客户端封装：初始化、ping、IK 验证、索引幂等创建 | `pkg/es/client.go` | 0.3 | 可连接 ES 并创建索引 |
| 3.2 | ES 同步服务：upsert / delete / bulk rebuild | `pkg/es/sync.go` | 3.1 | 单条同步和全量重建可用 |
| 3.3 | ES 搜索服务：multi_match（best_fields）+ 高亮 + create_time 降序 + 降级 PG LIKE | `pkg/es/search.go` | 3.1 | 搜索返回高亮，按时间降序；ES 不可用时降级 PG LIKE |
| 3.4 | 文章 service 集成 ES 同步 | `internal/service/article.go` | 2.6, 3.2 | 创建/更新/删除自动同步 ES |
| 3.5 | 文章 ES 搜索接口 + 重建索引接口 | `internal/handler/article.go` | 3.3, 3.4 | `/api/article/es/search` 和 `/api/article/es/rebuild` 可用 |
| 3.6 | OSS 客户端封装：定义 `ObjectStorage` 接口 + AliyunOSS 实现（Upload/Delete/TestConnectivity） | `pkg/oss/oss.go` | 0.3 | |
| 3.7 | OSS 配置运行时读取（从 blog_site_config） | `pkg/oss/oss.go` | 2.7 | 每次上传使用最新配置 |
| 3.8 | OSS 上传接口（文件命名：`blog/{year}/{month}/{uuid}.{ext}`） | `internal/handler/upload.go` | 3.6, 3.7 | 文件上传返回 OSS URL |
| 3.9 | OSS 连通测试接口（上传 4 字节测试文件后删除） | `internal/handler/site.go` | 3.6, 3.7 | 测试返回成功/失败 |

### Phase 4：路由注册与启动

| 编号 | 需求描述 | 涉及文件 | 依赖 | 验收标准 |
|------|----------|----------|------|----------|
| 4.1 | 完整路由注册：分组 + 中间件绑定（JWT/RBAC/限流） | `router/router.go` | Phase 1-3 全部 | 所有接口按 §5 注册 |
| 4.2 | main.go 完整启动流程 | `cmd/main.go` | 4.1 | 启动流程：config → PG → Redis → ES check → JWT init → router → ListenAndServe |
| 4.3 | Redis 缓存层完善：权限缓存、站点配置缓存集成 | 各 service 层 | 4.1 | 缓存命中/失效正常 |

### Phase 5：验证与文档

| 编号 | 需求描述 | 涉及文件 | 依赖 | 验收标准 |
|------|----------|----------|------|----------|
| 5.1 | 全接口联调测试 | 全部 | Phase 4 | 所有接口功能正确 |
| 5.2 | ES 降级场景验证 | 全部 | 5.1 | 关闭 ES 后搜索降级 PG LIKE |
| 5.3 | Redis 故障降级验证 | 全部 | 5.1 | Redis 不可用时系统不崩溃 |
| 5.4 | 编写 README：部署、初始化管理员、ES 重建索引 | `README.md` | 5.1 | 文档完整可操作 |

---

## 12. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| ES 服务宕机 | 文章搜索不可用 | 自动降级 PG LIKE 查询；日志告警；提供手动重建索引接口 |
| Redis 服务故障 | token 鉴权、权限缓存、限流失效 | 降级直查 PG；系统不崩溃 |
| OSS 配置错误 | 文件上传失败 | 保存配置时执行连通测试；上传失败返回明确错误信息 |
| ES IK 分词器未安装 | 索引创建失败或搜索异常 | 启动时 `_analyze` 验证；未安装时标记降级模式，日志 ERROR |
| ES 与 PG 数据不一致 | 搜索结果与实际数据不匹配 | 提供手动重建全量索引接口；同步失败记录日志 |
| 超级管理员误删 | 系统无法管理 | 代码层面硬编码禁止删除超级管理员角色和账号 |
| 密码泄露 | 账号安全 | bcrypt 加密；登录限流 5次/分钟/IP；logout 清除 Redis token |
| config.yaml 含敏感信息 | 安全风险 | 生产环境通过环境变量覆盖；OSS 配置存 PG 不写配置文件 |
| 并发角色修改导致缓存不一致 | 权限判断错误 | 角色修改时查询该角色下所有用户，批量清除权限缓存 |
| PG 连接池耗尽 | 接口超时/报错 | 合理配置 max_open_conns / max_idle_conns；连接超时监控 |
