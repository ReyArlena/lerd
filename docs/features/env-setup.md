# Environment Setup

`lerd env` sets up the `.env` file for a Laravel project in one command:

```bash
cd ~/Lerd/my-app
lerd env
```

---

## What it does

1. **Creates `.env`** from `.env.example` if no `.env` exists yet
2. **Detects which services the project uses** by inspecting the existing env keys — `DB_CONNECTION`, `REDIS_HOST`, `MAIL_HOST`, `SCOUT_DRIVER`, `FILESYSTEM_DISK`, `BROADCAST_CONNECTION`, etc.
3. **Writes lerd connection values** for each detected service (hosts, ports, credentials) — preserving all comments and line order
4. **Creates the project database** (and a `<name>_testing` database) inside the running service container; reports if they already exist
5. **Starts any referenced service** that is not already running
6. **Sets `APP_URL`** to the project's registered `.test` domain (`https://` if secured, `http://` otherwise)
7. **Generates `APP_KEY`** via `php artisan key:generate` if the key is missing or empty

---

## Example output

```
Creating .env from .env.example...
  Detected mysql        — applying lerd connection values
  Detected redis        — applying lerd connection values
  Detected mailpit      — applying lerd connection values
  Setting APP_URL=http://my-app.test
  Generating APP_KEY...
Done.
```

---

## Safe to re-run

Running `lerd env` on a project that already has a `.env` is safe — it only updates connection-related keys and leaves everything else untouched.
