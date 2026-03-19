# AI Integration (MCP)

Lerd ships a [Model Context Protocol](https://modelcontextprotocol.io/) server, letting AI assistants (Claude Code, JetBrains Junie, and any other MCP-compatible tool) manage your dev environment directly — run migrations, start services, toggle queue workers, and inspect logs without leaving the chat.

---

## Injecting the config

Run this once from a Laravel project root:

```bash
cd ~/Lerd/my-app
lerd mcp:inject
```

This writes three files:

| File | Purpose |
|---|---|
| `.mcp.json` | MCP server entry for Claude Code |
| `.claude/skills/lerd/SKILL.md` | Skill file that teaches Claude about lerd tools |
| `.junie/mcp/mcp.json` | MCP server entry for JetBrains Junie |

The command **merges** into existing configs — other MCP servers (e.g. `laravel-boost`, `herd`) are left untouched. Re-running it is safe.

To target a different directory:

```bash
lerd mcp:inject --path ~/Lerd/another-app
```

---

## Available MCP tools

Once the MCP server is connected, your AI assistant has access to:

| Tool | Description |
|---|---|
| `sites` | List all registered lerd sites (name, domain, path, PHP version, queue status) |
| `artisan` | Run `php artisan` in the PHP-FPM container — migrations, generators, seeders, cache, tinker |
| `service_start` | Start an infrastructure service (mysql, redis, postgres, …) |
| `service_stop` | Stop a service |
| `queue_start` | Start a queue worker for a site |
| `queue_stop` | Stop a queue worker |
| `logs` | Fetch recent container logs (nginx, any service, PHP version, or site name) |

---

## Example interactions

```
You: run migrations for the whitewaters project
AI:  → sites()           # finds path /home/user/Lerd/whitewaters
     → artisan(path: "/home/user/Lerd/whitewaters", args: ["migrate"])
     ✓  Ran 3 migrations in 42ms

You: the app is throwing 500s — check the logs
AI:  → logs(target: "8.4", lines: 50)
     PHP Fatal error: Class "App\Jobs\ProcessOrder" not found ...
```
