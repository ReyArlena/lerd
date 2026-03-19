# Lerd

**Laravel Herd for Linux** — a Podman-native local development environment for Laravel projects.

Lerd bundles Nginx, PHP-FPM, and optional services (MySQL, Redis, PostgreSQL, Meilisearch, MinIO) as rootless Podman containers, giving you automatic `.test` domain routing, per-project PHP/Node version isolation, and one-command TLS — all without touching your system's PHP or web server.

---

## Lerd vs Laravel Sail

[Laravel Sail](https://laravel.com/docs/sail) is the official per-project Docker Compose solution. Lerd is a shared infrastructure approach, closer to what [Laravel Herd](https://herd.laravel.com/) does on macOS. Both are valid — they solve slightly different problems.

|  | Lerd | Laravel Sail |
|---|---|---|
| Nginx | One shared container for all sites | Per-project |
| PHP-FPM | One container per PHP version, shared | Per-project container |
| Services (MySQL, Redis…) | One shared instance | Per-project (or manually shared) |
| `.test` domains | Automatic, zero config | Manual `/etc/hosts` or dnsmasq |
| HTTPS | `lerd secure` → trusted cert instantly | Manual or roll your own mkcert |
| RAM with 5 projects running | ~200 MB | ~1–2 GB (5× stacks) |
| Requires changes to project files | No | Yes — needs `docker-compose.yml` committed |
| Works on legacy / client repos | Yes — just `lerd link` | Only if you can add Sail |
| Defined in code (infra-as-code) | No | Yes |
| Team parity (all OS) | Linux only | macOS, Windows, Linux |

**Choose Sail when:** your team uses it, you need per-project service versions, or you want infrastructure defined in the repo.

**Choose Lerd when:** you work across many projects at once and don't want a separate stack per repo, you can't modify project files, you want instant `.test` routing, or you're on Linux and want the Herd experience.

---

## Next steps

- [Requirements](getting-started/requirements.md) — what you need before installing
- [Installation](getting-started/installation.md) — one-line installer or build from source
- [Quick Start](getting-started/quick-start.md) — up and running in 60 seconds
