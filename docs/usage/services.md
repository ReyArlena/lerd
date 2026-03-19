# Services

## Commands

| Command | Description |
|---|---|
| `lerd service start <name>` | Start a service (auto-installs on first use) |
| `lerd service stop <name>` | Stop a service container |
| `lerd service restart <name>` | Restart a service container |
| `lerd service status <name>` | Show systemd unit status |
| `lerd service list` | Show all services and their current state |

Available services: `mysql`, `redis`, `postgres`, `meilisearch`, `minio`, `mailpit`, `soketi`.

---

## Service credentials

!!! tip "Two sets of hostnames"
    Services run as Podman containers on the `lerd` network. Two hostnames apply depending on where you're connecting from:

    - **From host tools** (e.g. TablePlus, Redis CLI): use `127.0.0.1`
    - **From your Laravel app** (PHP-FPM runs inside the `lerd` network): use container hostnames (e.g. `lerd-mysql`)

    `lerd service start <name>` prints the correct `.env` variables to paste into your project.

| Service | Host (host tools) | Host (Laravel `.env`) | Port | User | Password | DB |
|---|---|---|---|---|---|---|
| MySQL | 127.0.0.1 | lerd-mysql | 3306 | root | `lerd` | `lerd` |
| PostgreSQL | 127.0.0.1 | lerd-postgres | 5432 | postgres | `lerd` | `lerd` |
| Redis | 127.0.0.1 | lerd-redis | 6379 | — | — | — |
| Meilisearch | 127.0.0.1 | lerd-meilisearch | 7700 | — | — | — |
| MinIO | 127.0.0.1 | lerd-minio | 9000 | `lerd` | `lerdpassword` | — |
| Mailpit SMTP | 127.0.0.1 | lerd-mailpit | 1025 | — | — | — |
| Soketi | 127.0.0.1 | lerd-soketi | 6001 | — | — | — |

Additional UIs:

- MinIO console: `http://127.0.0.1:9001`
- Mailpit web UI: `http://127.0.0.1:8025`
- Soketi metrics: `http://127.0.0.1:9601`
