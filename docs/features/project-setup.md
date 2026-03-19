# Project Setup

`lerd setup` automates the standard steps for getting a fresh Laravel clone running locally. Run it from the project root:

```bash
cd ~/Projects/my-app
lerd setup
```

A checkbox list appears with all available steps pre-selected based on the current project state. Toggle steps with space, confirm with enter, then watch them run sequentially.

```
→ Registering site...
Linked: my-app -> my-app.test (PHP 8.4, Node 24)

? Select setup steps to run:
  ◉ composer install
  ◉ npm ci
  ◉ lerd env
  ◯ lerd mcp:inject
  ◉ php artisan migrate
  ◯ php artisan db:seed
  ◉ npm run build
  ◯ lerd secure
  ◉ lerd open
```

**`lerd link` always runs first** (before the checkbox UI appears) to ensure the site is registered with the correct PHP version. This is mandatory — it handles both freshly cloned projects and projects already auto-registered by a parked directory.

---

## Smart defaults

| Step | Default | Condition |
|---|---|---|
| `composer install` | - [x] on | only if `vendor/` is missing |
| `npm ci` | - [x] on | only if `node_modules/` is missing and `package.json` exists |
| `lerd env` | - [x] on | always |
| `lerd mcp:inject` | - [ ] off | opt-in |
| `php artisan migrate` | - [x] on | always |
| `php artisan db:seed` | - [ ] off | opt-in |
| `npm run build` | - [x] on | only if `package.json` exists |
| `lerd secure` | - [ ] off | opt-in |
| `lerd open` | - [x] on | always |

The asset build step detects the right command from `package.json` — it looks for `build`, `production`, or `prod` scripts in priority order.

---

## Error handling

If a step fails, you are prompted to continue or abort:

```
✗ migrate failed: exit status 1
  Continue with remaining steps? [y/N]:
```

---

## Flags

| Flag | Description |
|---|---|
| `--all` / `-a` | Select all steps without showing the prompt (CI/automation) |
| `--skip-open` | Skip opening the browser at the end |
