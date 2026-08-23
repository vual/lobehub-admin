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

Note: the default compose file uses the published `calciumion/new-api:latest` image. To run the code from this repository, build a local image and change the `image` value of the `new-api` service in the compose file:

~~~bash
docker build -t lobehub-admin:local .
~~~

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

