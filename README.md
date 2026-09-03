# Blog Admin

**博客后台管理系统 — Gin + PostgreSQL + Redis + Elasticsearch + React 19**

基于 RBAC 权限体系的博客后台管理系统，包含完整的前后端：
- **后端**：Gin + PostgreSQL + Redis + Elasticsearch，提供 RESTful API
- **前端**：React 19 + TypeScript + Ant Design 6 + UnoCSS，管理后台界面

功能涵盖：仪表盘、文章管理（含 ES 全文检索）、分类管理、用户管理、角色与菜单权限管理、多云 OSS 文件上传、站点配置管理。

---

## 本次更新（Upgrade Highlights）

### 新增功能

- **仪表盘页面**：新增 `/dashboard` 路由，展示站点统计数据（文章 / 分类 / 用户 / 运行天数）和系统运行环境信息（应用版本、Go 版本、内存占用、依赖健康状态）
- **菜单管理 CRUD**：菜单管理页面从只读升级为完整增删改，支持目录 / 页面 / 按钮三种类型
- **动态路由**：前端支持基于菜单数据的动态路由注册，新增页面只需在组件注册表（`componentRegistry`）注册 + 数据库建菜单即可生效
- **多云 OSS**：对象存储支持阿里云、腾讯云、七牛云、Amazon S3 四种厂商切换
- **常用网站图标**：常用链接支持设置 icon 图标
- **前端路由权限守卫**：基于菜单树的 `RequirePermission` 组件，非授权路由自动重定向

### 技术升级

- **代码分割**：使用 `React.lazy` + 组件注册表实现页面按需加载
- **4K 适配**：内容区限宽 1600px / 1800px，侧边栏大屏加宽，MDEditor 高度自适应
- **登录页优化**：入场动画、记住用户名、表单样式优化
- **分类时间格式化**：创建时间显示为 `YYYY:MM:DD:HH:MM:SS`

### 新增迁移文件

| 文件 | 说明 |
|------|------|
| `012_seed_menu.sql` | 菜单种子数据 + 编辑演示角色 |
| `013_add_oss_provider.sql` | OSS 厂商和区域字段 |
| `014_seed_dashboard_menu.sql` | 仪表盘菜单 |
| `015_add_link_icon.sql` | 常用链接图标字段 |

### 新增文件

| 文件 | 说明 |
|------|------|
| `pkg/storage/factory.go` | 多云存储工厂（按厂商配置创建存储客户端） |
| `pkg/storage/tencent.go` / `qiniu.go` / `s3.go` | 腾讯云 / 七牛云 / Amazon S3 存储实现 |
| `internal/service/system.go` / `internal/handler/system.go` | 系统运行环境信息接口（仪表盘） |
| `web/src/router/componentRegistry.ts` | 前端组件注册表（动态路由 + 懒加载） |
| `web/src/pages/Dashboard/index.tsx` | 仪表盘页面 |
| `internal/dto/menu.go` | 菜单 DTO（增删改请求/响应结构体） |

---

## 技术栈

### 后端

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| 后端框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 | 高性能 HTTP 框架 |
| ORM | [GORM](https://gorm.io) v1.31 | 数据库 ORM |
| 数据库 | PostgreSQL 16 | 主数据存储 |
| 缓存 | Redis（[go-redis/v8](https://github.com/go-redis/redis)） | Token 缓存、权限缓存、限流 |
| 全文检索 | Elasticsearch 8.x + IK 分词器 | 文章全文搜索 |
| 对象存储 | 多云支持：阿里云 OSS / 腾讯云 COS / 七牛云 Kodo / Amazon S3 | 图片/文件上传，厂商可切换 |
| 认证 | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | JWT Token 鉴权 |
| 配置管理 | [Viper](https://github.com/spf13/viper) v1.21 | YAML 配置加载 |
| 验证码 | [go-captcha/v2](https://github.com/wenlng/go-captcha) | 滑动验证码 |
| API 文档 | [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) | Swagger 在线文档 |

### 前端

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| 框架 | [React 19](https://react.dev) | 声明式 UI 框架 |
| 语言 | [TypeScript](https://www.typescriptlang.org) | 类型安全的 JavaScript 超集 |
| UI 组件库 | [Ant Design 6](https://ant.design) | 企业级 React UI 组件库 |
| CSS 方案 | [UnoCSS](https://unocss.dev) | 原子化 CSS 引擎（自定义 4K 断点） |
| 路由 | [React Router](https://reactrouter.com) v6 | 路由管理（含动态路由） |
| 状态管理 | [Zustand](https://zustand.docs.pmnd.rs) | 轻量状态管理（认证/权限） |
| HTTP 客户端 | [Axios](https://axios-http.com) | HTTP 请求库 |
| 构建工具 | [Vite](https://vitejs.dev) | 下一代前端构建工具 |

---

## 项目结构

```
gin-blog-server/
├── cmd/
│   └── main.go                    # 应用入口，启动流程编排
├── config/
│   ├── config.go                  # 配置结构体定义与加载
│   └── config.yaml                # 配置文件（数据库/Redis/ES/JWT 等）
├── internal/                      # 内部业务逻辑
│   ├── dto/                       # 数据传输对象（含 menu.go 菜单 DTO）
│   ├── handler/                   # HTTP 处理器（Controller 层，含 system.go）
│   ├── model/                     # 数据模型（GORM 模型定义）
│   ├── repository/                # 数据访问层（数据库操作）
│   └── service/                   # 业务逻辑层（含 system.go 系统概览）
├── middleware/                    # 中间件（JWT 鉴权、RBAC 权限、限流）
├── migrations/                    # 数据库迁移 SQL（001~015，按序执行）
├── pkg/                           # 公共工具包
│   ├── captcha/                   # 滑动验证码
│   ├── elasticsearch/             # ES 客户端封装（索引管理/同步/搜索）
│   ├── errors/                    # 业务错误码定义
│   ├── jwt/                       # JWT 工具（生成/解析 Token）
│   ├── logger/                    # 日志工具
│   ├── migration/                 # 迁移文件执行器
│   ├── redis/                     # Redis 客户端封装
│   ├── response/                  # 统一 JSON 响应格式
│   ├── storage/                   # 多云对象存储（factory 工厂 + oss/tencent/qiniu/s3 实现）
│   └── util/                      # 通用工具（分页等）
├── router/
│   └── router.go                  # 路由注册（分组 + 中间件绑定）
├── docs/                          # Swagger 自动生成文档
├── web/                           # 管理后台前端（React 19 + TS + Ant Design + UnoCSS）
│   ├── src/
│   │   ├── api/                   # API 接口封装（按模块分文件）
│   │   ├── assets/                # 静态资源
│   │   ├── components/            # 公共组件（Layout 布局、PageHeader 等）
│   │   ├── router/                # 路由配置（index.tsx 动态路由 + componentRegistry.ts 组件注册表）
│   │   ├── store/                 # Zustand 状态管理（认证信息、站点状态）
│   │   ├── pages/                 # 页面组件（Dashboard/Article/Category/User/Role/Menu/Site/...）
│   │   ├── App.tsx                # 根组件
│   │   └── main.tsx               # 入口文件
│   ├── unocss.config.ts           # UnoCSS 配置（自定义 4K 断点）
│   ├── vite.config.ts             # Vite 配置（代理、别名等）
│   └── package.json
├── blog/                          # 博客展示端前端（访客浏览界面）
├── go.mod
├── go.sum
├── PRD.md                         # 产品需求文档
└── REQUIREMENTS.md                # 需求清单与分阶段实现计划
```

---

## 环境依赖

运行本项目前，请确保以下环境已安装：

### 后端依赖

| 依赖 | 版本要求 | 说明 |
|------|----------|------|
| Go | 1.26+ | 后端编译运行环境 |
| PostgreSQL | 16 | 主数据库 |
| Redis | 6+ | 缓存与限流 |
| Elasticsearch | 8.x | 全文检索引擎 |
| IK 分词器 | 与 ES 版本一致 | ES 中文分词插件 |

### 前端依赖

| 依赖 | 版本要求 | 说明 |
|------|----------|------|
| Node.js | 18+ | 前端运行环境 |
| npm | 9+ | 前端包管理工具 |

### IK 分词器安装

ES 必须预装 IK 分词器插件，否则文章全文搜索将无法正常工作：

```bash
./bin/elasticsearch-plugin install https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-8.17.0.zip
```

> **注意**：IK 分词器版本号必须与 Elasticsearch 版本一致。安装后需重启 ES。

---

## 快速开始

### 1. 克隆项目

```bash
git clone <repo-url>
cd gin-blog-server
```

### 2. 后端配置

编辑 `config/config.yaml`，根据实际环境修改以下配置：

```yaml
database:
  host: 127.0.0.1        # PostgreSQL 地址
  port: 5432
  user: postgres          # 数据库用户名
  password: your_password # 数据库密码
  dbname: blog            # 数据库名称

redis:
  addr: 127.0.0.1:6379   # Redis 地址
  password: ""
  db: 1

elasticsearch:
  addresses:
    - "http://127.0.0.1:9200"  # ES 地址
  index_prefix: "blog"
  username: ""
  password: ""

jwt:
  secret: your-jwt-secret-key  # JWT 签名密钥（请修改为随机字符串）
  expire: 24                   # Token 过期时间（小时）
  prefix: blog_                # Token 前缀标识
```

> OSS 配置不在 yaml 中，运行时通过管理后台页面配置（支持阿里云 / 腾讯云 / 七牛云 / Amazon S3 切换），存储于数据库 `blog_site_config` 表。

### 3. 数据库初始化

`migrations/` 目录包含 001~015 编号迁移文件（建表 + 初始数据），应用启动时会按 `migration.dir` 配置自动按序执行，也可手动导入：

```bash
for f in migrations/*.sql; do psql -U postgres -d blog < "$f"; done
```

迁移文件概览：

| 文件 | 说明 |
|------|------|
| `001_init.sql` | 基础建表（用户/角色/菜单/分类/文章/站点配置）+ 初始数据 |
| `002_add_profile.sql` | 站长个人资料 |
| `003_add_site_founded.sql` | 站点创建时间字段 |
| `004_seed_data.sql` / `005_seed_site_config.sql` | 种子数据 |
| `006_add_icp.sql` / `007_add_police_icp.sql` | ICP 备案字段 |
| `008_article_covers.sql` | 文章封面 |
| `009_add_blog_profile.sql` / `010_add_blog_link.sql` | 博客资料与常用链接 |
| `011_add_repost_fields.sql` | 文章转载字段 |
| `012_seed_menu.sql` | 菜单种子数据 + 编辑演示角色 |
| `013_add_oss_provider.sql` | OSS 厂商和区域字段 |
| `014_seed_dashboard_menu.sql` | 仪表盘菜单 |
| `015_add_link_icon.sql` | 常用链接图标字段 |

### 4. 启动后端服务

```bash
go mod tidy
go run cmd/main.go
```

后端服务启动后将按以下流程初始化：配置加载 → PostgreSQL 连接 → Redis 连接 → ES 连通性检查 → 路由注册 → 监听端口（默认 `:3000`）。

### 5. 启动前端服务

```bash
cd web
npm install
npm run dev
```

前端开发服务器默认运行在 `http://localhost:5173`，已配置代理将 `/api` 请求转发至后端服务。

> 博客展示端位于 `blog/` 目录，启动方式相同。

---

## API 文档

统一响应格式：

```json
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

> `code=200` 成功；`401` 未登录；`403` 无权限；其他非 200 值为业务错误码。

> 在线文档：启动后端后访问 `/swagger/index.html`（需配置 `swagger.enabled: true`）。

### 接口概览

#### 认证模块

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 账号密码登录，返回 JWT Token |
| POST | `/api/auth/logout` | 登出，清除 Redis Token |
| GET | `/api/auth/userinfo` | 获取当前用户信息 + 菜单树 + 权限标识 |

#### RBAC 权限

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/menu/list` | 获取全部菜单树 |
| POST | `/api/menu/create` | 新增菜单（目录/页面/按钮） |
| PUT | `/api/menu/update` | 更新菜单 |
| DELETE | `/api/menu/:id` | 删除菜单 |
| GET | `/api/role/list` | 角色列表 |
| POST | `/api/role/create` | 创建角色 |
| PUT | `/api/role/update` | 更新角色与权限 |
| DELETE | `/api/role/:id` | 删除角色 |

#### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/user/page` | 用户分页列表 |
| POST | `/api/user/create` | 创建用户 |
| PUT | `/api/user/update` | 编辑用户 |
| PUT | `/api/user/resetPwd/:id` | 重置密码 |
| DELETE | `/api/user/:id` | 删除用户 |

#### 分类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/category/tree` | 获取分类树 |
| POST | `/api/category/create` | 新增分类 |
| PUT | `/api/category/update` | 更新分类 |
| DELETE | `/api/category/:id` | 删除分类（被文章引用时拦截） |

#### 文章管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/article/page` | PG 普通分页列表 |
| GET | `/api/article/es/search` | ES 全文搜索（支持高亮） |
| POST | `/api/article/create` | 创建文章，自动同步 ES |
| PUT | `/api/article/update` | 更新文章，自动同步 ES |
| DELETE | `/api/article/:id` | 逻辑删除，同步删除 ES 文档 |
| POST | `/api/article/es/rebuild` | 手动重建 ES 全量索引（仅超级管理员） |

#### 仪表盘 & 系统

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/site/stats` | 站点统计（文章/分类/用户/运行天数） |
| GET | `/api/system/overview` | 系统运行概览（版本/Go 版本/内存/依赖健康状态） |

#### 站点 & OSS

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/site/config` | 获取站点配置 |
| PUT | `/api/site/config` | 保存配置，刷新 Redis 缓存 |
| POST | `/api/site/testOss` | OSS 连通性测试 |
| POST | `/api/upload/oss` | form-data 文件上传，返回 OSS URL |
| GET / POST | `/api/site/links` | 常用链接列表 / 新增（支持 icon） |
| PUT / DELETE | `/api/site/links/:id` | 更新 / 删除常用链接 |

> 详细请求/响应参数请参见 [PRD.md](./PRD.md) 第 6 节及 Swagger 在线文档。

---

## Redis 缓存说明

| Key 模式 | 类型 | TTL | 用途 |
|----------|------|-----|------|
| `token:{userId}` | string | 24h | 登录 Token 存储，登出时主动删除 |
| `rbac:perm:{userId}` | string (JSON) | 24h | 用户权限标识集合缓存 |
| `rbac:menu:{userId}` | string (JSON) | 24h | 用户可访问菜单树缓存 |
| `site:config` | string (JSON) | 1h | 站点全局配置缓存 |
| `rate:login:{ip}` | string (int) | 60s | 登录接口限流计数（5 次/分钟/IP） |

**缓存失效触发点：**

- **角色修改**（菜单/权限变更）→ 批量删除该角色下所有用户的 `rbac:perm` 和 `rbac:menu`
- **用户信息变更**（编辑/禁用/删除）→ 删除该用户的 `token`、`rbac:perm`、`rbac:menu`
- **站点配置保存** → 删除 `site:config`

---

## Elasticsearch 说明

### 索引配置

- **索引名**：`blog_article_index`
- **分词器**：`ik_max_word`（索引时最大切词）、`ik_smart`（搜索时智能分词）
- **索引字段**：`id`、`title`、`content`、`tags`、`category_id`、`status`、`create_time`

### 同步策略

| 触发场景 | ES 操作 |
|----------|---------|
| 文章创建 | upsert 文档到 ES |
| 文章更新 | upsert 文档到 ES |
| 文章逻辑删除 | 从 ES 删除对应文档 |
| 手动重建索引 | 全量从 PG 读取未删除文章，bulk upsert |

### 降级机制

ES 不可用时，全文搜索自动降级为 PostgreSQL 模糊查询：

```sql
WHERE is_deleted = 0 AND (title LIKE '%keyword%' OR content LIKE '%keyword%' OR tags LIKE '%keyword%')
```

应用启动时会检查 ES 连通性和 IK 分词器可用性，检查失败不阻止启动，仅日志告警并进入降级模式。

### 重建索引

当 ES 与 PG 数据不一致时，可使用超级管理员账号调用重建接口：

```
POST /api/article/es/rebuild
```

该接口将删除现有索引、重新创建并全量同步 PG 中所有未删除文章。

---

## 初始化管理员

系统首次部署后，需通过 SQL 手动插入超级管理员账号：

```sql
-- 1. 插入超级管理员角色
INSERT INTO sys_role (role_name, menu_ids, button_perms, is_super, create_time, update_time)
VALUES ('超级管理员', '[]', '[]', 1, NOW(), NOW());

-- 2. 插入超级管理员用户（密码为 bcrypt 加密后的值，此处示例密码为 Admin@123）
INSERT INTO sys_user (username, password, nickname, role_id, status, create_time, update_time)
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '超级管理员', 1, 1, NOW(), NOW());
```

> **注意**：上述密码为 `Admin@123` 的 bcrypt 哈希值。生产环境请务必修改密码。可使用以下方式生成新的 bcrypt 密码哈希：
>
> ```bash
> # 使用 htpasswd 工具
> htpasswd -nbBC 10 "" "your-password" | tr -d ':\n' | sed 's/$2y/$2a/'
> ```

菜单种子数据由 `migrations/001_init.sql`、`migrations/012_seed_menu.sql`、`migrations/014_seed_dashboard_menu.sql` 自动导入。

---

## 相关文档

- [PRD.md](./PRD.md) — 产品需求文档（数据模型、API 详细定义、业务规则）
- [REQUIREMENTS.md](./REQUIREMENTS.md) — 需求清单与分阶段实现计划
