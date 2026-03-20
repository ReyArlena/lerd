# Queue Workers

Lerd can run Laravel queue workers as persistent systemd user services. The worker runs `php artisan queue:work` inside the project's PHP-FPM container and restarts automatically on failure.

## Commands

| Command | Description |
|---|---|
| `lerd queue:start` | Start a queue worker for the current project |
| `lerd queue:stop` | Stop the queue worker for the current project |
| `lerd queue start` | Same as `queue:start` (subcommand form) |
| `lerd queue stop` | Same as `queue:stop` (subcommand form) |

---

## Options for `queue:start`

| Flag | Default | Description |
|---|---|---|
| `--queue` | `default` | Queue name to process |
| `--tries` | `3` | Max attempts before marking a job as failed |
| `--timeout` | `60` | Seconds a job may run before timing out |

---

## Redis requirement

If `QUEUE_CONNECTION=redis` is set in the project's `.env`, lerd verifies that `lerd-redis` is running before starting the worker. If it is not, you will see:

```
queue worker requires Redis (QUEUE_CONNECTION=redis in .env) but lerd-redis is not running
Start it first: lerd services start redis
```

---

## Example

```bash
cd ~/Lerd/my-app
lerd queue:start --queue=emails,default --tries=5 --timeout=120
# Systemd unit: lerd-queue-my-app.service
# Logs: journalctl --user -u lerd-queue-my-app -f
```

---

## Web UI control

Queue workers are controllable from the **Sites tab** in the web UI. The amber toggle starts or stops the worker. When a worker is running, a **Queue** log tab appears in the site detail panel alongside PHP-FPM. The amber dot next to the site in the sidebar indicates the worker is active.
