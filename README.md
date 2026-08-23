# LobeHub Admin

基于 [new-api](https://github.com/QuantumNous/new-api) 构建的 AI API 网关与 LobeHub 管理平台。

本项目保留 new-api 的完整能力，并增加 LobeHub 用户管理功能，方便在同一个后台完成 API 网关运维和 LobeHub 账号治理。

## 主要新增功能

### LobeHub 用户管理

登录管理员后台后，可在侧边栏的「LobeHub → 用户管理」中：

- 分页查看并按 ID、用户名、邮箱、姓名或手机号搜索用户
- 按启用状态、封禁状态、角色、邮箱验证和 2FA 状态筛选
- 查看用户资料、登录提供商、会话数量、最后活跃时间和认证状态
- 编辑用户名、邮箱、头像、手机号、姓名及验证状态
- 封禁或解封用户，支持填写封禁原因和设置过期时间
- 将用户角色设为 `user` 或 `admin`；覆盖自定义角色时需要显式确认
- 重置用户密码并生成一次性临时密码
- 变更角色、封禁或重置密码时清理可撤销的数据库会话和相关凭据
- 对管理员操作记录审计信息，便于追踪敏感变更

## LobeHub 集成要求

LobeHub 用户管理读取的是已有 LobeHub PostgreSQL 业务表，不会替 LobeHub 创建或升级表结构。

需要满足以下条件：

1. 网关主数据库必须是 PostgreSQL；网关本身仍可使用 SQLite 或 MySQL，但此时无法使用 LobeHub 用户管理。
2. LobeHub 业务表必须位于同一个 PostgreSQL 数据库中，并且网关连接账号具备读取和写入权限。
3. `LOBEHUB_DB_SCHEMA` 配置为 LobeHub 业务表所在 schema，默认为 `public`。
4. LobeHub 数据库结构需要包含 `users`、`accounts`、`auth_sessions` 以及 OIDC 相关表和字段。

示例配置：

~~~dotenv
SQL_DSN=postgresql://user:password@127.0.0.1:5432/lobehub
LOBEHUB_DB_SCHEMA=public
REDIS_CONN_STRING=redis://127.0.0.1:6379/0
SESSION_SECRET=请替换为足够长的随机字符串
~~~

如果 LobeHub 表位于独立 schema，例如 `lobehub`：

~~~dotenv
LOBEHUB_DB_SCHEMA=lobehub
~~~

启动后如果 schema 不存在或字段不兼容，管理页会提示对应错误；请检查 schema 配置和 LobeHub 数据库版本，不要直接对生产库执行未经验证的迁移。

## 快速开始

### 使用 Docker Compose

仓库中的 `docker-compose.yml` 适合快速启动 PostgreSQL、Redis 和网关服务。首次使用前，请修改其中的数据库密码、Redis 密码和会话密钥：

~~~bash
git clone <your-repository-url>
cd lobehub-admin

# 修改 docker-compose.yml 中的密码和环境变量
docker compose up -d
~~~

访问：<http://localhost:3000>

注意：默认 compose 文件使用已发布的 `calciumion/new-api:latest` 镜像。如果需要运行当前仓库的代码，请先构建本地镜像，并将 compose 中 `new-api` 服务的 `image` 改为本地镜像：

~~~bash
docker build -t lobehub-admin:local .
~~~

### 本地开发

环境要求：

- Go 1.22 及以上
- Bun
- PostgreSQL 15 或兼容版本
- Redis

构建前端并启动后端：

~~~bash
cd web
bun install --frozen-lockfile
bun run build

cd ..
go run main.go
~~~

## 常用配置

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `PORT` | HTTP 服务端口 | `3000` |
| `SQL_DSN` | 主数据库连接字符串 | - |
| `SQL_SCHEMA` | PostgreSQL 网关业务 schema | `admin` |
| `LOG_SQL_DSN` | 独立日志数据库连接字符串 | - |
| `REDIS_CONN_STRING` | Redis 连接字符串 | - |
| `LOBEHUB_DB_SCHEMA` | LobeHub 业务表所在 PostgreSQL schema | `public` |
| `SESSION_SECRET` | 登录会话签名密钥，多节点必须一致 | - |
| `SESSION_COOKIE_SECURE` | HTTPS 部署时启用 Secure Cookie | `false` |
| `SESSION_COOKIE_TRUSTED_URL` | Secure Cookie 模式下允许的 HTTPS Origin | - |
| `TRUSTED_PROXIES` | 可信反向代理 IP/CIDR | - |
| `STREAMING_TIMEOUT` | 流式请求无响应超时时间，单位为秒 | `300` |
| `MEMORY_CACHE_ENABLED` | 是否启用内存缓存 | - |

更多变量及说明请查看 [.env.example](./.env.example)。

## 管理员首次登录

首次启动后访问首页，按初始化向导创建管理员账号。完成初始化后：

1. 登录管理后台。
2. 在系统设置中配置渠道、模型、计费和安全策略。
3. 确认 `LOBEHUB_DB_SCHEMA` 指向正确的 LobeHub schema。
4. 从侧边栏进入「LobeHub → 用户管理」验证用户列表和操作权限。

LobeHub 用户管理入口只对网关管理员开放；修改 LobeHub 全局角色属于高风险操作，需要更高权限，并可能使用户现有登录凭据失效。

## 项目结构

~~~text
.
├── controller/       HTTP 请求处理
├── model/            数据模型与数据库访问
├── service/          业务逻辑
├── relay/            AI 上游适配与请求转发
├── router/           API、管理端和网页路由
├── middleware/       认证、限流、日志和安全中间件
├── web/              React + TypeScript 管理前端
├── relaykit/         可独立构建的 Relay 工具模块
└── docs/             API、部署和认证相关文档
~~~

LobeHub 用户管理相关实现主要位于：

- `controller/lobehub_user.go`
- `service/lobehub_user.go`
- `model/lobehub_user.go`
- `web/src/features/lobehub-users/`
- `web/src/routes/_authenticated/lobehub/users/`

## 开发与检查

~~~bash
# 前端类型检查、构建和测试
cd web
bun run typecheck
bun run build
bun run test

# 回到项目根目录执行 Go 测试
cd ..
GOWORK=off go test ./...

# relaykit 必须独立构建
cd relaykit
GOWORK=off go build ./...
~~~

提交涉及 LobeHub 用户管理的变更时，建议同时验证：权限校验、并发更新、会话撤销、PostgreSQL schema 配置和前端多语言显示。

## 相关文档

- [环境变量示例](./.env.example)
- [认证、会话与 PAT 说明](./docs/authentication.md)
- [OpenAPI 文档](./docs/openapi/api.json)
- [Relay API 文档](./docs/openapi/relay.json)
- [宝塔部署说明](./docs/installation/BT.md)
- [new-api 官方文档](https://docs.newapi.pro/en/docs)
- [new-api DeepWiki](https://deepwiki.com/QuantumNous/new-api)

## 免责声明

本项目仅用于合法、合规且获得授权的 AI API 聚合、内部管理、模型调用和用户服务场景。使用者需要自行获取上游 API 密钥、账号及模型服务授权，并遵守上游服务条款、所在地法律法规以及适用的数据保护和内容安全要求。

## 交流
微信：822784588

