# LobeHub Admin

[English](./README.en.md)

基于 [new-api](https://github.com/QuantumNous/new-api) 构建的 AI API 网关与 [LobeHub](https://github.com/lobehub/lobehub) 管理平台。

本项目保留 new-api 的完整能力，并增加 LobeHub 管理功能，方便在同一个后台完成 API 网关运维和 LobeHub 后台管理。

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

### LobeHub 会话管理

登录超级管理员后台后，可在侧边栏的「LobeHub → 会话管理」中：

- 分页查看 LobeHub 会话，并按会话、用户、Agent 或群组进行搜索
- 按会话类型、状态、触发方式、模型、提供商和更新时间筛选
- 按创建时间、更新时间、消息数、Token 总量或总费用排序
- 查看会话所属用户、Agent/群组、消息数量、模型、提供商和更新时间
- 打开会话详情，按时间顺序查看消息内容、角色、推理、工具调用、搜索信息、翻译、附件、用量和元数据
- 使用游标分页加载较早消息；该功能只读，不会修改或删除 LobeHub 会话及消息

会话管理入口仅对超级管理员开放。会话内容可能包含用户输入、模型输出、附件链接和业务元数据，请根据组织的数据访问和隐私规范配置管理员权限。

### LobeHub 知识库 / RAG 管理

登录管理员后台后，可在侧边栏的「LobeHub → 知识库」中：

- 分页查看全部个人知识库和 Workspace 知识库，并按 ID、名称、创建者、Workspace、可见性和 RAG 状态搜索或筛选
- 查看文件、文档、切片数量、向量覆盖率、存储大小、创建者和 Workspace 映射
- 查看文件级分块 / 向量任务状态、任务错误、文档解析正文、页面数据、编辑器数据和切片文本
- 仅修改知识库名称、描述和头像；使用更新时间进行乐观锁校验，避免覆盖其他管理员的修改
- 只展示向量是否存在和向量模型，不读取或返回向量数组

知识库管理入口：`/lobehub/knowledge-bases`。该模块只读取和更新 LobeHub PostgreSQL 表，不修改 LobeHub 的接口、代码或表结构。

## LobeHub 集成要求

LobeHub 用户管理读取的是已有 LobeHub PostgreSQL 业务表，不会替 LobeHub 创建或升级表结构。

需要满足以下条件：

1. 网关主数据库必须是 PostgreSQL；网关本身仍可使用 SQLite 或 MySQL，但此时无法使用 LobeHub 用户管理。
2. LobeHub 业务表必须位于同一个 PostgreSQL 数据库中，并且网关连接账号具备读取和写入权限。
3. `LOBEHUB_DB_SCHEMA` 配置为 LobeHub 业务表所在 schema，默认为 `public`。
4. LobeHub 用户管理需要 `users`、`accounts`、`auth_sessions` 以及 OIDC 相关表和字段；启用会话管理时，还需要 `topics`、`messages`、`threads`、`message_groups`、`agents`、`chat_groups` 及消息附件、插件、翻译、搜索查询和 TTS 相关表及字段。
5. 启用知识库管理时，还需要 `knowledge_bases`、`knowledge_base_files`、`files`、`documents`、`async_tasks`、`file_chunks`、`chunks` 和 `embeddings` 表及本模块读取的字段。管理端会在访问时校验 schema、表和字段；缺失或不兼容时只返回 `LOBEHUB_SCHEMA_UNAVAILABLE`。

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

默认 compose 文件使用上游 `calciumion/new-api:latest` 镜像。如果要运行本项目的 LobeHub Admin 镜像，请将 compose 中 `new-api` 服务的 `image` 改为 `registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest`，或使用 `ghcr.io/vual/lobehub-admin:latest`：

~~~bash
docker pull registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

### 与 LobeHub 一起启动

如果你使用 LobeHub 项目中的 `docker-compose/deploy/docker-compose.yml`，可以将本项目作为同一个 Compose 项目的服务加入。这样管理端会与 LobeHub 共享 `lobe-network`、PostgreSQL 和 Redis，并直接读取 LobeHub 的用户表。

在 LobeHub Compose 文件的 `services:` 下追加以下服务。默认直接拉取阿里云镜像，也可以将 `image` 替换为 GHCR 镜像：

~~~yaml
  lobehub-admin:
    image: registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
    container_name: lobehub-admin
    restart: always
    command: ["--log-dir", "/app/logs"]
    ports:
      - "3000:3000"
    volumes:
      - ./lobehub-admin-data:/data
      - ./lobehub-admin-logs:/app/logs
    environment:
      - SQL_DSN=postgresql://postgres:${POSTGRES_PASSWORD}@postgresql:5432/${LOBE_DB_NAME}
      - LOBEHUB_DB_SCHEMA=public
      - REDIS_CONN_STRING=redis://redis:6379/0
      - SESSION_SECRET=replace-with-a-long-random-secret
      - SESSION_COOKIE_SECURE=false
      - TZ=Asia/Shanghai
    depends_on:
      postgresql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - lobe-network
~~~

然后在 LobeHub Compose 目录执行：

~~~bash
docker compose pull
docker compose up -d
docker compose ps
~~~

启动后：

- LobeHub：<http://localhost:3210>
- LobeHub Admin：<http://localhost:3000>

容器之间必须使用 Compose 服务名通信，因此 `SQL_DSN` 使用 `postgresql`，Redis 使用 `redis`，不要在容器配置中写 `localhost`。LobeHub Compose 默认将 PostgreSQL 的业务表放在 `public` schema；如果实际使用其他 schema，请同步修改 `LOBEHUB_DB_SCHEMA`。

`SESSION_COOKIE_SECURE=false` 仅适用于本地 HTTP 调试。生产环境通过 HTTPS 访问时，应设置为 `true`，并同时配置 `SESSION_COOKIE_TRUSTED_URL`。请为 `SESSION_SECRET` 设置高强度随机值，不要使用示例值。

### 独立启动本项目 Docker

如果不需要与 LobeHub 使用同一个 Compose 项目，可以使用本仓库的独立 Compose 配置：

1. 拉取已发布镜像：

~~~bash
docker pull registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

2. 修改本仓库 `docker-compose.yml` 中 `new-api` 服务的镜像：

~~~yaml
services:
  new-api:
    image: registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

3. 启动独立服务：

~~~bash
docker compose pull
docker compose up -d
docker compose ps
~~~

该配置会独立启动网关、PostgreSQL 和 Redis，网关地址为 <http://localhost:3000>。如果还要启用 LobeHub 用户管理，`SQL_DSN` 必须连接到包含 LobeHub 业务表的 PostgreSQL 数据库，并设置 `LOBEHUB_DB_SCHEMA`；不能使用独立 Compose 新建的空数据库代替 LobeHub 数据库。

也可以只运行本项目容器，并连接外部 PostgreSQL 和 Redis：

~~~bash
docker run -d --name lobehub-admin --restart unless-stopped \
  -p 3000:3000 \
  -e SQL_DSN="postgresql://postgres:password@host.docker.internal:5432/lobechat" \
  -e LOBEHUB_DB_SCHEMA=public \
  -e REDIS_CONN_STRING="redis://host.docker.internal:6379/0" \
  -e SESSION_SECRET="replace-with-a-long-random-secret" \
  -e SESSION_COOKIE_SECURE=false \
  -e TZ=Asia/Shanghai \
  -v ./data:/data \
  -v ./logs:/app/logs \
  registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

`host.docker.internal` 仅适用于 Docker Desktop 常见场景；Linux 或远程数据库请替换为实际可访问的主机名或 IP。生产环境请使用独立的密钥管理、HTTPS、数据库备份和最小权限数据库账号。


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
4. 从侧边栏进入「LobeHub → 用户管理」验证用户列表和操作权限；超级管理员还可以进入「LobeHub → 会话管理」查看会话和消息详情，管理员可以进入「LobeHub → 知识库」查看知识库和 RAG 状态。

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
- `controller/lobehub_conversation.go`
- `service/lobehub_conversation.go`
- `model/lobehub_conversation.go`
- `controller/lobehub_knowledge_base.go`
- `service/lobehub_knowledge_base.go`
- `model/lobehub_knowledge_base.go`
- `web/src/features/lobehub-users/`
- `web/src/routes/_authenticated/lobehub/users/`
- `web/src/features/lobehub-conversations/`
- `web/src/routes/_authenticated/lobehub/conversations/`
- `web/src/features/lobehub-knowledge-bases/`
- `web/src/routes/_authenticated/lobehub/knowledge-bases/`

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

提交涉及 LobeHub 用户或会话管理的变更时，建议同时验证：权限校验、并发更新、会话撤销、会话详情分页、PostgreSQL schema 配置和前端多语言显示。

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
