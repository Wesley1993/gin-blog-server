# PRD‑Blog‑Admin v0.0.1（博客后台管理系统）

> 
> AI‑Coding 专用 PRD，可直接丢给 AI 生成代码
> 技术栈：**Gin + PostgreSQL16 + Redis + Elasticsearch + Vue3+TS+Element‑Plus**
> 权限模型：**标准 RBAC (用户‑角色‑菜单 / 权限)**；ES 用于**文章全文检索**；Redis 做 token 缓存、权限缓存、限流
> 仅管理后台，**不含博客前台**

## 1 项目概述

**项目名称**：Blog 后台管理系统
**版本**：v1.0
**目标**：博客后台，RBAC 权限体系、用户 / 角色 / 菜单权限管理；文章、分类、站点配置；OSS 文件上传；**ES 实现文章全文搜索**；Redis 缓存 token 与权限数据。
**使用场景**：Web 浏览器后台；使用者：超级管理员、编辑角色。

**非目标（暂时不处理，后期可能继续扩展）**

1. 不开发博客前台页面
2. 不做第三方登录，仅账号密码登录
3. 文件不落地本地，全部上传 OSS
4. 不实现评论、访问统计模块
5. ES 只负责文章检索，不承担业务数据库，业务数据以 PG 为准

## 2 功能需求 FR

### FR‑01 RBAC 权限管理（用户‑角色‑菜单‑权限）

- 功能描述：完整 RBAC 模型；角色绑定菜单与按钮权限；用户绑定角色；接口层做权限校验；权限数据放入 Redis 缓存，减少 PG 查询。
- 输入：角色名称、菜单 ID 列表、权限标识；用户绑定角色 ID。
- 输出：角色列表、权限树；登录后返回当前用户可访问菜单、权限标识集合。
- 业务规则：
  1. 内置超级管理员角色，不可删除，拥有全部权限。
  2. 用户绑定**单个角色**；一个角色可绑定多个菜单、多个按钮权限标识。
  3. 用户登录后，角色权限树存入 Redis 缓存；角色修改后主动清除对应 Redis 权限缓存。
  4. 无权限菜单前端隐藏；无权限按钮隐藏；接口鉴权不通过返回`403`。
  5. 菜单分为目录、页面、按钮；按钮绑定权限标识如`article:add`、`article:edit`。

### FR‑02 用户管理

- 功能描述：后台账号增删改查、启用禁用、重置密码。
- 输入：用户名、昵称、角色 ID、账号状态、密码。
- 输出：用户分页列表、表单。
- 业务规则：
  1. 用户名 PG 唯一约束；密码 bcrypt 加密存储。
  2. 账号禁用后禁止登录；超级管理员账号禁止删除。
  3. 用户信息变更，清理该用户 Redis 内 token、权限缓存。

### FR‑03 分类管理

- 功能描述：文章分类树形管理，排序、启用禁用。
- 输入：分类名称、父分类 ID、排序值、状态。
- 输出：树形分类列表。
- 业务规则：
  1. 支持一 / 二级分类；sort 越小越靠前。
  2. 被文章引用的分类允许禁用，**禁止删除**。

### FR‑04 文章管理（对接 ES 全文索引）

- 功能描述：文章增删改查，富文本编辑，标签；草稿 / 发布状态；OSS 图片上传；**新增 / 更新 / 删除文章时同步 ES 索引；后台支持 ES 全文检索文章**。
- 输入：标题、分类 ID、标签、富文本 content、cover 封面、status 状态。
- 输出：PG 分页列表；ES 全文搜索结果列表。
- 业务规则：
  1. 文章状态：`0草稿`，`1已发布`；逻辑删除`is_deleted=1`。
  2. 富文本图片上传 OSS，返回 URL 存入正文。
  3. **PG 为真值源；文章创建 / 更新 / 逻辑删除，触发 ES 文档同步**。
  4. ES 索引字段：id、title、content、tags、category_id、status、create_time。
  5. ES 搜索：支持标题、正文、标签模糊全文检索；仅查询 ES，详情从 PG 读取。
  6. 草稿文章也写入 ES，后台可搜到；**ES 不同步逻辑删除的数据**。

### FR‑05 站点管理 & OSS 配置

- 功能描述：站点基础配置、OSS 对象存储配置。
- 输入：站点名称、描述、版权；OSS AK/SK/Bucket/Endpoint/ 自定义域名。
- 输出：配置表单，保存生效。
- 业务规则：
  1. OSS 配置存入 PG；上传接口读取该配置；保存时执行 OSS 连通性测试。
  2. 站点配置写入 Redis 缓存，减少 DB 查询。

## 3 界面交互定义（Vue3 后台）

```
整体布局：侧边菜单 + 顶部导航 + 主内容区
页面清单：
1.登录页：账号密码登录
2.权限管理：角色列表、角色编辑弹窗、权限树多选；菜单管理页面
3.用户管理：分页表格，新增/编辑/重置密码/禁用
4.分类管理：树形表格，新增子分类、排序
5.文章管理：
    -普通分页列表；
    -全文搜索框：调用ES搜索；
    -新增/编辑页面，富文本编辑器，OSS图片上传
6.站点设置：Tab【基础设置】【OSS存储设置】

交互行为：
1.全部列表支持分页、筛选；文章页增加全文搜索输入框。
2.保存loading；toast提示；删除二次确认。
3.OSS保存执行连通测试，返回成功/失败。
4.富文本图片直接上传OSS，url插入编辑器。
5.权限变更后，当前用户需要重新登录刷新权限。
```

## 4 数据模型 PostgreSQL Schema

> 
> PG16；主键`bigserial`；json 字段使用`jsonb`；时间`timestamp without time zone`

> **逻辑删除说明**：所有表使用应用层管理 `is_deleted` 字段（`smallint NOT NULL DEFAULT 0`，0=未删 1=已删）做逻辑删除，**不使用 GORM 自动软删除**。所有查询需手动加 `WHERE is_deleted = 0`。

### 公共基础模型

所有表共有字段：

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 PRIMARY KEY |
| create_time | timestamp without time zone | Y | 创建时间 |
| update_time | timestamp without time zone | Y | 更新时间 |

> 各表定义中如已包含上述字段则不再重复列出 `id`，但 `create_time`、`update_time` 均保留在表定义中以便阅读。逻辑删除字段 `is_deleted` 仅在需要的表中出现。

### sys_user 用户表

表格

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 |
| username | varchar(50) | Y | 登录账号，UNIQUE 唯一 |
| password | varchar(100) | Y | bcrypt 加密密码 |
| nickname | varchar(50) | N | 昵称 |
| role_id | bigint | Y | 绑定角色 ID |
| status | smallint | Y | 0 禁用 1 启用 |
| create_time | timestamp without time zone | Y |  |
| update_time | timestamp without time zone | Y |  |

### sys_role 角色表

表格

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 |
| role_name | varchar(50) | Y | 角色名称 |
| menu_ids | jsonb | N | 菜单 id 数组 jsonb |
| button_perms | jsonb | N | 按钮权限标识数组，如 `["article:add","article:edit"]` |
| is_super | smallint | Y | 是否超级管理员 0/1 |
| create_time | timestamp without time zone | Y |  |
| update_time | timestamp without time zone | Y |  |

### sys_menu 菜单权限表

表格

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 |
| parent_id | bigint | Y | 父菜单 ID |
| menu_name | varchar(50) | Y | 菜单名称 |
| menu_type | smallint | Y | 1 目录 2 页面 3 按钮 |
| path | varchar(100) | N | 前端路由 path |
| perms | varchar(100) | N | 权限标识 article:add |
| sort | int | Y | 排序值 |
| status | smallint | Y | 0 禁用 1 启用 |
| create_time | timestamp without time zone | Y |  |
| update_time | timestamp without time zone | Y |  |

### blog_category 分类表

表格

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 |
| name | varchar(60) | Y | 分类名称 |
| parent_id | bigint | Y | 父分类 ID |
| sort | int | Y | 排序 |
| status | smallint | Y | 0 禁用 1 启用 |
| create_time | timestamp without time zone | Y |  |
| update_time | timestamp without time zone | Y |  |

### blog_article 文章表

表格

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 |
| title | varchar(120) | Y | 标题 |
| category_id | bigint | Y | 分类 ID |
| cover | varchar(255) | N | 封面 OSS 地址 |
| content | text | Y | 富文本正文 |
| tags | varchar(200) | N | 标签逗号分隔 |
| status | smallint | Y | 0 草稿 1 已发布 |
| is_deleted | smallint | Y | 逻辑删除 0 未删 1 已删 |
| create_time | timestamp without time zone | Y |  |
| update_time | timestamp without time zone | Y |  |

**索引定义**：

```sql
CREATE INDEX idx_article_category ON blog_article(category_id);
CREATE INDEX idx_article_status ON blog_article(status);
CREATE INDEX idx_article_create_time ON blog_article(create_time);
```

### blog_site_config 站点配置表

表格

| 字段 | PG 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | bigserial | Y | 主键 |
| site_name | varchar(100) | N | 站点名称 |
| site_desc | varchar(255) | N | 站点描述 |
| copyright | varchar(255) | N | 版权 |
| oss_access_key | varchar(200) | N |  |
| oss_secret_key | varchar(200) | N |  |
| oss_bucket | varchar(100) | N |  |
| oss_endpoint | varchar(200) | N |  |
| oss_domain | varchar(200) | N |  |
| create_time | timestamp without time zone | Y |  |
| update_time | timestamp without time zone | Y |  |

### ES index: `blog_article_index`

> 
> ES 仅做检索，数据来源 PG，不做写入源

```
{
  "mappings": {
    "properties": {
      "id": {"type":"long"},
      "title": {"type":"text","analyzer":"ik_max_word","search_analyzer":"ik_smart"},
      "content": {"type":"text","analyzer":"ik_max_word","search_analyzer":"ik_smart"},
      "tags": {"type":"keyword"},
      "category_id": {"type":"long"},
      "status": {"type":"integer"},
      "create_time": {"type":"date"}
    }
  }
}
```

同步策略：

1. 文章新建 / 更新：upsert 到 ES
2. 文章逻辑删除：delete 对应 ES 文档
3. 提供手动重建索引接口，全量从 PG 同步到 ES

### ES 实现细节

- **IK 分词器依赖**：ES 必须预装 Ik 插件。安装命令示例：
  ```bash
  ./bin/elasticsearch-plugin install https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-8.17.0.zip
  ```
  （版本号需与 ES 版本一致）
- **启动时验证**：应用启动时调用 `GET /_analyze` `{"analyzer":"ik_max_word","text":"测试中文分词"}` 验证 IK 可用性。不可用则标记 ES 为降级模式，日志打印 ERROR，**不阻止应用启动**。
- **索引管理**：应用启动时幂等检查并创建索引（先 `Indices.Exists`，不存在则 `Indices.Create`）。
- **ES 客户端**：使用 `github.com/elastic/go-elasticsearch/v8`。
- **search_analyzer**：mapping 中 title/content 字段增加 `"search_analyzer": "ik_smart"`，索引用 `ik_max_word` 最大切词提高召回，搜索用 `ik_smart` 智能分词提高精度。
- **降级行为**：ES 不可用时，全文搜索降级为 PG 查询：
  ```sql
  WHERE is_deleted=0 AND (title LIKE '%keyword%' OR content LIKE '%keyword%' OR tags LIKE '%keyword%')
  ```
  日志打印 ES 异常。

## 5 Redis 使用说明

1. `token:{userId}`：存储登录 token，设置过期时间，登录鉴权。**TTL 24h**（与 JWT 过期时间一致）。
2. `rbac:perm:{userId}`：缓存用户权限标识集合，角色变更主动删除 key。**TTL 24h**。
3. `rbac:menu:{userId}`：缓存用户可访问菜单树。**TTL 24h**。
4. `site:config`：缓存站点全局配置。**TTL 1h**。
5. 接口限流：`rate:login:{ip}` — 限制 5次/分钟/IP，使用 Redis `INCR + EXPIRE` 实现。

### 缓存失效触发点

- **角色修改**（菜单/权限变更）→ 查询该角色下所有用户 → 批量删除其 `rbac:perm:{userId}` 和 `rbac:menu:{userId}`
- **用户信息变更**（编辑/禁用/删除）→ 删除该用户 `token:{userId}`、`rbac:perm:{userId}`、`rbac:menu:{userId}`
- **站点配置保存** → 删除 `site:config` key

## 5.1 OSS 实现细节

- **SDK 选型**：`github.com/aliyun/alibabacloud-oss-go-sdk-v2`（官方推荐，V4 签名）。
- **连通测试方式**：上传极小测试文件（4字节）后立即删除，同时验证读写权限。
- **文件命名规则**：`blog/{year}/{month}/{uuid}.{ext}`，按日期分目录，UUID 防冲突。
- **统一接口**：定义 `ObjectStorage` 接口（Upload / Delete / TestConnectivity），当前实现 AliyunOSS，未来可扩展 S3。
- **OSS 配置来源**：运行时从 PG `blog_site_config` 表动态读取，不从 config.yaml 读取。

## 6 API 接口定义

统一返回格式

```
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

> 
> code=200 成功；401 未登录；403 无权限；非 0 业务错误码

**分页参数规范**：

- `page`：页码，从 1 开始，默认 1
- `page_size`：每页条数，默认 10，最大 100
- 分页响应格式：
```json
{"code":200,"msg":"ok","data":{"list":[...],"total":100,"page":1,"page_size":10}}
```

### 错误码表

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 401 | 未登录/token过期 |
| 403 | 无权限 |
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

### 认证接口

- POST `/api/auth/login` 账号密码登录，返回 token
  ```json
  // 请求
  {"username":"admin","password":"123456"}
  // 响应
  {"code":200,"msg":"ok","data":{"token":"eyJhbGci..."}}
  ```
- POST `/api/auth/logout` 登出，清除 redis token
- GET `/api/auth/userinfo` 返回当前登录用户信息 + 可访问菜单树 + 权限标识集合
  ```json
  // 响应
  {"code":200,"msg":"ok","data":{"user":{"id":1,"username":"admin","nickname":"管理员","role_id":1},"menus":[{"id":1,"parent_id":0,"menu_name":"权限管理","menu_type":1,"path":"/system","children":[...]}],"permissions":["article:add","article:edit","article:delete"]}}
  ```

### RBAC 权限

1. GET `/api/menu/list` 获取全部菜单树
   ```json
   // 响应
   {"code":200,"msg":"ok","data":[{"id":1,"parent_id":0,"menu_name":"权限管理","menu_type":1,"path":"/system","sort":1,"status":1,"children":[...]}]}
   ```
2. GET `/api/role/list` 角色列表
   ```json
   // 响应
   {"code":200,"msg":"ok","data":[{"id":1,"role_name":"超级管理员","menu_ids":[1,2,3],"button_perms":["article:add"],"is_super":1}]}
   ```
3. POST `/api/role/create` 创建角色
   ```json
   // 请求
   {"role_name":"编辑","menu_ids":[1,2,3],"button_perms":["article:add","article:edit"]}
   ```
4. PUT `/api/role/update` 更新角色 + 权限
   ```json
   // 请求
   {"id":2,"role_name":"编辑","menu_ids":[1,2,3],"button_perms":["article:add","article:edit"]}
   ```
5. DELETE `/api/role/:id` 删除角色

### 用户管理

1. GET `/api/user/page` 用户分页
   ```json
   // 请求参数: ?page=1&page_size=10
   // 响应
   {"code":200,"msg":"ok","data":{"list":[{"id":1,"username":"admin","nickname":"管理员","role_id":1,"status":1}],"total":50,"page":1,"page_size":10}}
   ```
2. POST `/api/user/create` 创建用户
   ```json
   // 请求
   {"username":"editor","nickname":"编辑","password":"123456","role_id":2,"status":1}
   ```
3. PUT `/api/user/update` 编辑用户
   ```json
   // 请求
   {"id":2,"nickname":"编辑A","role_id":2,"status":1}
   ```
4. PUT `/api/user/resetPwd/:id` 重置密码
   ```json
   // 请求
   {"password":"newpassword"}
   ```
5. DELETE `/api/user/:id` 删除用户

### 分类管理

1. GET `/api/category/tree` 获取分类树
   ```json
   // 响应
   {"code":200,"msg":"ok","data":[{"id":1,"name":"技术","parent_id":0,"sort":1,"status":1,"children":[{"id":2,"name":"Go","parent_id":1,"sort":1,"status":1}]}]}
   ```
2. POST `/api/category/create` 新增
   ```json
   // 请求
   {"name":"Go","parent_id":1,"sort":1,"status":1}
   ```
3. PUT `/api/category/update` 更新
   ```json
   // 请求
   {"id":2,"name":"Golang","parent_id":1,"sort":1,"status":1}
   ```
4. DELETE `/api/category/:id` 删除（被引用拦截）

### 文章管理（含 ES 搜索）

1. GET `/api/article/page` PG 普通分页列表
   ```json
   // 请求参数: ?page=1&page_size=10&category_id=1&status=1
   // 响应
   {"code":200,"msg":"ok","data":{"list":[{"id":1,"title":"Hello Go","category_id":1,"cover":"https://oss.../cover.jpg","tags":"go,web","status":1,"create_time":"2025-01-01 00:00:00"}],"total":100,"page":1,"page_size":10}}
   ```
2. GET `/api/article/es/search` ES 全文搜索接口，query 参数 keyword
   ```json
   // 请求参数: ?keyword=Go%E5%BC%80%E5%8F%91&page=1&page_size=10
   // 响应
   {"code":200,"msg":"ok","data":{"list":[{"id":1,"title":"Hello Go","tags":"go,web","category_id":1,"status":1,"create_time":"2025-01-01 00:00:00"}],"total":5,"page":1,"page_size":10}}
   ```
3. POST `/api/article/create` 创建文章，内部同步 ES
   ```json
   // 请求
   {"title":"Hello Go","category_id":1,"cover":"https://oss.../cover.jpg","content":"<p>正文内容</p>","tags":"go,web","status":1}
   ```
4. PUT `/api/article/update` 更新文章，内部同步 ES
   ```json
   // 请求
   {"id":1,"title":"Hello Go v2","category_id":1,"content":"<p>更新后正文</p>","tags":"go","status":1}
   ```
5. DELETE `/api/article/:id` 逻辑删除，同步删除 ES 文档
6. POST `/api/article/es/rebuild` 手动重建 ES 全量索引（仅超级管理员）

### 站点 & OSS

1. GET `/api/site/config` 获取配置
   ```json
   // 响应
   {"code":200,"msg":"ok","data":{"site_name":"我的博客","site_desc":"一个技术博客","copyright":"©2025","oss_access_key":"***","oss_secret_key":"***","oss_bucket":"my-blog","oss_endpoint":"oss-cn-hangzhou.aliyuncs.com","oss_domain":"cdn.example.com"}}
   ```
2. PUT `/api/site/config` 保存配置，刷新 redis 缓存
   ```json
   // 请求
   {"site_name":"我的博客","site_desc":"一个技术博客","copyright":"©2025","oss_access_key":"LTAI...","oss_secret_key":"xxx","oss_bucket":"my-blog","oss_endpoint":"oss-cn-hangzhou.aliyuncs.com","oss_domain":"cdn.example.com"}
   ```
3. POST `/api/site/testOss` OSS 连通测试
   ```json
   // 响应
   {"code":200,"msg":"ok","data":null}
   ```
4. POST `/api/upload/oss` form‑data 文件上传，返回 url
   ```json
   // 请求: form-data file=@image.png
   // 响应
   {"code":200,"msg":"ok","data":{"url":"https://cdn.example.com/blog/2025/01/a1b2c3d4.png"}}
   ```

## 7 非功能需求 NFR

- 技术栈
  - 前端：Vue3 + TS + Element‑Plus + Vue‑Router + Pinia
  - 后端：Gin + GORM + PostgreSQL16 + Redis + Elasticsearch8.x
  - OSS：兼容阿里云 OSS/S3 协议对象存储
  - ES 分词：ik 分词器
- 鉴权：token+RBAC 权限标识校验；权限、站点配置使用 Redis 缓存。
- ES 容错：ES 服务不可用时，**降级走 PG 数据库模糊查询**，不直接崩溃。
- 配置说明：
  - `config.yaml` 不包含 rabbitmq、order 配置（博客系统无此需求）
  - OSS 参数由数据库动态管理，不在 yaml 中配置
  - JWT token prefix 为 `blog_`
- 输出产物：
  1. 完整前后端可运行代码
  2. PG 建表 SQL
  3. ES index mapping 定义
  4. Redis key 说明文档
  5. README：部署、初始化管理员、ES 重建索引操作

### 中间件

- **JWT 鉴权中间件** (`middleware/jwt.go`)：从 `Authorization: Bearer {token}` 提取 token → 校验 JWT 有效性 → 校验 Redis 中 `token:{userId}` 是否存在（防止 token 被主动注销）→ 将 `userID`、`roleID`、`isSuper` 注入 `gin.Context`。
- **RBAC 权限中间件** (`middleware/rbac.go`)：从 Redis 读取 `rbac:perm:{userId}` 权限标识集合 → 获取当前路由所需的权限标识（通过路由 metadata 或自定义 header）→ 比对，不匹配返回 403 → 超级管理员（`is_super=1`）跳过权限检查。
- **限流中间件**：登录接口专用，基于 IP 的 5次/分钟限流。

## 8 异常处理清单

1. token 过期返回 401，前端跳转登录。
2. 接口无权限返回 403。
3. Redis 故障：降级直查 PG，系统不崩溃。
4. ES 不可用：全文搜索降级为 PG 数据库 like 模糊查询，日志打印 ES 异常。
5. OSS 参数错误，保存时报错提示。
6. 删除被引用分类，返回业务提示禁止删除。
7. 角色修改，自动清除对应用户权限缓存。
8. ES 同步失败记录日志，提供重建索引接口修复。

## 9 验收标准

- RBAC：超级管理员拥有全部权限；普通角色只能看到分配菜单与按钮；接口层权限拦截生效。
- 用户增删改查、禁用账号、重置密码；超级管理员不可删除。
- 分类树形展示，排序、禁用，被引用分类禁止删除。
- 文章新增 / 更新 / 删除自动同步 ES；全文搜索可命中标题、正文、标签；ES 宕机可降级 PG 模糊搜索。
- OSS 配置保存、连通测试正常；图片上传返回 OSS 地址。
- Redis 缓存权限与站点配置；角色变更缓存失效。
- 手动重建 ES 索引接口正常，可修复 ES 数据不一致。