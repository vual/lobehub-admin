# LobeHub Admin

An AI API gateway and LobeHub administration platform based on [new-api](https://github.com/QuantumNous/new-api).

This project retains the complete capabilities of new-api and adds LobeHub user management, allowing AI gateway operations and LobeHub account administration to be handled from one console.

## Main Addition

### LobeHub User Management

After signing in to the administrator console, open **LobeHub → Users** in the sidebar to:

- Browse users with pagination and search by ID, username, email, name, or phone number
- Filter users by status, ban status, role, email verification, and 2FA status
- View user profiles, login providers, session counts, last active time, and authentication status
- Edit usernames, email addresses, avatars, phone numbers, names, and verification states
- Ban or unban users with a reason and optional expiration time
- Set a user's role to `user` or `admin`; explicit confirmation is required before overwriting a custom role
- Reset user passwords and generate one-time temporary passwords
- Revoke revocable database sessions and related credentials when changing roles, banning users, or resetting passwords
- Record administrator actions for auditing sensitive changes

## LobeHub Integration Requirements

LobeHub user management reads the existing LobeHub PostgreSQL business tables. It does not create or upgrade the LobeHub schema.

The following requirements must be met:

1. The gateway's primary database must be PostgreSQL. The gateway itself still supports SQLite and MySQL, but LobeHub user management is unavailable with those databases.
2. LobeHub business tables must be in the same PostgreSQL database, and the gateway database user must have read and write permissions.
3. Set `LOBEHUB_DB_SCHEMA` to the schema containing the LobeHub business tables. The default is `public`.
4. The LobeHub database must contain `users`, `accounts`, `auth_sessions`, and the required OIDC-related tables and columns.

Example configuration:

~~~dotenv
SQL_DSN=postgresql://user:password@127.0.0.1:5432/lobehub
LOBEHUB_DB_SCHEMA=public
REDIS_CONN_STRING=redis://127.0.0.1:6379/0
SESSION_SECRET=replace-with-a-long-random-secret
~~~

If the LobeHub tables are in a separate schema such as `lobehub`:

~~~dotenv
LOBEHUB_DB_SCHEMA=lobehub
~~~

If the schema does not exist or its columns are incompatible, the administration page will report an error. Check the schema configuration and LobeHub database version instead of applying unverified migrations to a production database.

## Quick Start

### Using Docker Compose

The repository's `docker-compose.yml` can quickly start PostgreSQL, Redis, and the gateway. Before starting, change the database password, Redis password, and session secret:

~~~bash
git clone <your-repository-url>
cd lobehub-admin

# Edit passwords and environment variables in docker-compose.yml
docker compose up -d
~~~

Open <http://localhost:3000>.

The default compose file uses the upstream `calciumion/new-api:latest` image. To run the LobeHub Admin image for this project, change the `image` value of the `new-api` service to `registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest`, or use `ghcr.io/vual/lobehub-admin:latest`:

~~~bash
docker pull registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

### Run Together with LobeHub

If you use LobeHub's `docker-compose/deploy/docker-compose.yml`, you can add this project as another service in the same Compose project. The administration service will then share LobeHub's `lobe-network`, PostgreSQL, and Redis, and can read LobeHub's user tables directly.

Append the following service under `services:` in the LobeHub Compose file. The service pulls the Aliyun image by default; you can replace `image` with the GHCR image if needed:

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

Then run the following commands from the LobeHub Compose directory:

~~~bash
docker compose pull
docker compose up -d
docker compose ps
~~~

After startup:

- LobeHub: <http://localhost:3210>
- LobeHub Admin: <http://localhost:3000>

Containers must communicate through Compose service names. Therefore, use `postgresql` in `SQL_DSN` and `redis` for Redis; do not use `localhost` inside container configuration. The LobeHub Compose setup places PostgreSQL business tables in the `public` schema by default. If you use another schema, update `LOBEHUB_DB_SCHEMA` accordingly.

`SESSION_COOKIE_SECURE=false` is suitable only for local HTTP development. For production HTTPS access, set it to `true` and configure `SESSION_COOKIE_TRUSTED_URL` as well. Set `SESSION_SECRET` to a strong random value and never use the example value in production.

### Run This Project with Standalone Docker

If you do not need to use the same Compose project as LobeHub, use this repository's standalone Compose configuration:

1. Pull the published image:

~~~bash
docker pull registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

2. Change the image for the `new-api` service in this repository's `docker-compose.yml`:

~~~yaml
services:
  new-api:
    image: registry.cn-hangzhou.aliyuncs.com/ann-chat/lobehub-admin:latest
~~~

3. Start the standalone services:

~~~bash
docker compose pull
docker compose up -d
docker compose ps
~~~

This configuration starts the gateway, PostgreSQL, and Redis independently. The gateway is available at <http://localhost:3000>. To enable LobeHub user management as well, `SQL_DSN` must point to a PostgreSQL database containing the LobeHub business tables, and `LOBEHUB_DB_SCHEMA` must be set. The empty database created by the standalone Compose file cannot replace the LobeHub database.

You can also run only this project's container and connect it to external PostgreSQL and Redis services:

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

`host.docker.internal` is common with Docker Desktop. On Linux or when using a remote database, replace it with an actually reachable hostname or IP address. For production, use dedicated secret management, HTTPS, database backups, and a least-privilege database account.


### Local Development

Requirements:

- Go 1.22 or later
- Bun
- PostgreSQL 15 or a compatible version
- Redis

Build the frontend and start the backend:

~~~bash
cd web
bun install --frozen-lockfile
bun run build

cd ..
go run main.go
~~~

## Common Configuration

| Variable | Description | Default |
| --- | --- | --- |
| `PORT` | HTTP service port | `3000` |
| `SQL_DSN` | Primary database connection string | - |
| `SQL_SCHEMA` | PostgreSQL schema for gateway tables | `admin` |
| `LOG_SQL_DSN` | Separate log database connection string | - |
| `REDIS_CONN_STRING` | Redis connection string | - |
| `LOBEHUB_DB_SCHEMA` | PostgreSQL schema containing LobeHub tables | `public` |
| `SESSION_SECRET` | Login session signing secret; must be identical across nodes | - |
| `SESSION_COOKIE_SECURE` | Enable Secure Cookies for HTTPS deployments | `false` |
| `SESSION_COOKIE_TRUSTED_URL` | Allowed HTTPS Origins in Secure Cookie mode | - |
| `TRUSTED_PROXIES` | Trusted reverse-proxy IPs/CIDRs | - |
| `STREAMING_TIMEOUT` | Streaming no-response timeout in seconds | `300` |
| `MEMORY_CACHE_ENABLED` | Enable the in-memory cache | - |

See [.env.example](./.env.example) for the complete list of environment variables.

## First Administrator Login

After the first startup, open the homepage and follow the setup wizard to create an administrator account. Then:

1. Sign in to the administration console.
2. Configure channels, models, billing, and security policies in system settings.
3. Confirm that `LOBEHUB_DB_SCHEMA` points to the correct LobeHub schema.
4. Open **LobeHub → Users** in the sidebar to verify the user list and permissions.

The LobeHub user management section is available only to gateway administrators. Changing a LobeHub user's global role is a high-risk operation that requires elevated permissions and may invalidate the user's existing login credentials.

## Project Structure

~~~text
.
├── controller/       HTTP request handlers
├── model/            Data models and database access
├── service/          Business logic
├── relay/            AI upstream adapters and request forwarding
├── router/           API, administration, and web routes
├── middleware/       Authentication, rate limiting, logging, and security middleware
├── web/              React + TypeScript administration frontend
├── relaykit/         Independently buildable Relay utility module
└── docs/             API, deployment, and authentication documentation
~~~

LobeHub user management is mainly implemented in:

- `controller/lobehub_user.go`
- `service/lobehub_user.go`
- `model/lobehub_user.go`
- `web/src/features/lobehub-users/`
- `web/src/routes/_authenticated/lobehub/users/`

## Development Checks

~~~bash
# Frontend type checking, build, and tests
cd web
bun run typecheck
bun run build
bun run test

# Run Go tests from the project root
cd ..
GOWORK=off go test ./...

# relaykit must build independently
cd relaykit
GOWORK=off go build ./...
~~~

When changing LobeHub user management, also verify authorization, concurrent updates, session revocation, PostgreSQL schema configuration, and frontend localization.

## Documentation

- [Environment variable example](./.env.example)
- [Authentication, sessions, and PATs](./docs/authentication.md)
- [OpenAPI documentation](./docs/openapi/api.json)
- [Relay API documentation](./docs/openapi/relay.json)
- [BaoTa deployment guide](./docs/installation/BT.md)
- [new-api official documentation](https://docs.newapi.pro/en/docs)
- [new-api DeepWiki](https://deepwiki.com/QuantumNous/new-api)

## Disclaimer

This project is intended only for lawful, authorized AI API aggregation, internal administration, model access, and user-service scenarios. Users are responsible for obtaining upstream API keys, accounts, and model-service authorization, and for complying with upstream terms of service, applicable laws and regulations, data-protection requirements, and content-safety obligations.

## Contact

WeChat: 822784588
