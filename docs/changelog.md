# Changelog

All notable changes to Lerd will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Lerd uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.6.0] — 2026-03-21

### Added

- **Git worktree support** — each `git worktree` checkout automatically gets its own subdomain (`<branch>.<site>.test`) with a dedicated nginx vhost. No manual steps required.
  - The watcher daemon detects `git worktree add` / `git worktree remove` in real time via fsnotify and generates or removes vhosts accordingly. It watches `.git/` itself so it correctly re-attaches when `.git/worktrees/` is deleted (last worktree removed) and re-created (new worktree added).
  - Startup scan generates vhosts for all existing worktrees across all registered sites.
  - `EnsureWorktreeDeps` — symlinks `vendor/` and `node_modules/` from the main repo into each worktree checkout, and copies `.env` with `APP_URL` rewritten to the worktree subdomain.
  - `lerd sites` shows worktrees indented under their parent site.
  - The web UI shows worktrees in the site detail panel with clickable domain links and an open-in-browser button.
  - A git-branch icon appears on the site button in the sidebar whenever the site has active worktrees.
- **HTTPS for worktrees** — when a site is secured with `lerd secure`, all its worktrees automatically receive an SSL vhost that reuses the parent site's wildcard mkcert certificate (`*.domain.test`). No separate certificate is needed per worktree. Securing and unsecuring a site also updates `APP_URL` in each worktree's `.env`.
- **Catch-all default vhost** (`_default.conf`) — any `.test` hostname that does not match a registered site returns HTTP 444 / rejects the TLS handshake, instead of falling through to the first alphabetical vhost.
- **`stripe:listen` as a background service** — `lerd stripe:listen` now runs the Stripe CLI in a persistent systemd user service (`lerd-stripe-<site>.service`) rather than a foreground process. It survives terminal sessions and restarts on failure. `lerd stripe:listen stop` tears it down.
- **Service pause state** — `lerd service stop` now records the service as manually paused. `lerd start` and autostart on login skip paused services. `lerd stop` + `lerd start` restore the previous state: running services restart, manually stopped services stay stopped.
- **Queue worker Redis pre-flight** — `lerd queue:start` checks that `lerd-redis` is running when `QUEUE_CONNECTION=redis` is set in `.env`, and returns a friendly error with instructions rather than failing with a cryptic DNS error from PHP.

### Fixed

- **Park watcher depth** — the filesystem watcher no longer registers projects found in subdirectories of parked directories. Only direct children of a parked directory are eligible for auto-registration.
- **Nginx reload ordering for secure/unsecure** — `lerd secure` / `lerd unsecure` (and their UI/MCP equivalents) now save the updated `secured` flag to `sites.yaml` *before* reloading nginx. Previously a failed nginx reload would leave `sites.yaml` with a stale `secured` state, causing the watcher to regenerate the wrong vhost type on restart.
- **Tray always restarts on `lerd start`** — any existing tray process is killed before relaunching, preventing duplicate tray instances after repeated `lerd start` calls.
- **FPM quadlet skip-write optimisation** — `WriteFPMQuadlet` skips writing and daemon-reloading when the quadlet content is unchanged. Unnecessary daemon-reloads caused Podman's quadlet generator to regenerate all service files, which could briefly disrupt `lerd-dns` and cause `.test` resolution failures.

---

## [0.5.16] — 2026-03-20

### Fixed

- **PHP-FPM image build on restricted Podman** — fully qualify all base image names in the Containerfile (`docker.io/library/composer:latest`, `docker.io/library/php:X.Y-fpm-alpine`). Systems without unqualified-search registries configured in `/etc/containers/registries.conf` would fail with "short-name did not resolve to an alias".

---

## [0.5.15] — 2026-03-20

### Fixed

- **PHP-FPM image build on Podman** — the Containerfile now declares `FROM composer:latest AS composer-bin` as an explicit stage before copying the composer binary. Podman (unlike Docker) does not auto-pull images referenced only in `COPY --from`, causing builds to fail with "no stage or image found with that name". This also affected `lerd update` and `lerd php:rebuild` in v0.5.14, leaving containers stopped if the build failed after the old image was removed.
- **Zero-downtime PHP-FPM rebuild** — `lerd php:rebuild` no longer removes the existing image before building. The running container stays up during the build; only the final `systemctl restart` causes a brief interruption. Force rebuilds now use `--no-cache` instead of `rmi -f`.
- **UI logs panel** — clicking logs for a site whose PHP-FPM container is not running now shows a clean "container is not running" message instead of the raw podman error.
- **`lerd php` / `lerd artisan`** — running these when the PHP-FPM container is stopped now returns a friendly error with the `systemctl --user start` command instead of a raw podman error.
- **`lerd update` ensures PHP-FPM is running** — after applying infrastructure changes, `lerd update` now starts any installed PHP-FPM containers that are not running. Also fixed a cosmetic bug where "skipping rebuild" was printed even when a rebuild had just run.

---

## [0.5.14] — 2026-03-20

### Added

- **`LERD_SITE_PATH` in MCP config** — `mcp:inject` now embeds the project path as `LERD_SITE_PATH` in the injected MCP server config. The MCP server reads this at startup and uses it as the default `path` for `artisan`, `composer`, `env_setup`, `db_export`, and `site_link`, so AI assistants no longer need to pass an explicit path on every call.
- **`.ai/mcp/mcp.json` injection** — `mcp:inject` now also writes into `.ai/mcp/mcp.json` (used by Windsurf and other MCP-compatible tools), in addition to `.mcp.json` and `.junie/mcp/mcp.json`.

---

## [0.5.13] — 2026-03-20

### Fixed

- **`lerd mcp:inject`** — `db_export` tool now correctly appears in the generated SKILL.md and `.junie/guidelines.md` (skill content was omitted from the v0.5.12 release)

---

## [0.5.12] — 2026-03-20

### Added

- **MCP: `db_export`** — export the project database to a SQL dump file (defaults to `<database>.sql` in the project root); reads connection details from `.env`

### Fixed

- **`lerd artisan` / `lerd php` / `lerd node` / `lerd npm` / `lerd npx`** — lerd usage/help text and "Error: exit status N" no longer appear when the subprocess exits with a non-zero code (e.g. failed tests); only the subprocess output is shown and the original exit code is propagated to the shell

---

## [0.5.11] — 2026-03-20

### Added

- **MCP: 14 new tools** — the `lerd mcp` server now exposes the full project lifecycle:
  - `composer` — run Composer inside the PHP-FPM container
  - `node_install` / `node_uninstall` — install or uninstall Node.js versions via fnm
  - `runtime_versions` — list installed PHP and Node.js versions with defaults
  - `env_setup` — configure `.env` for lerd (detects services, starts them, creates DB, generates `APP_KEY`, sets `APP_URL`)
  - `site_link` / `site_unlink` — register or unregister a directory as a lerd site
  - `secure` / `unsecure` — enable or disable HTTPS for a site; updates `APP_URL` automatically
  - `xdebug_on` / `xdebug_off` / `xdebug_status` — toggle Xdebug per PHP version and check state
  - `service_add` / `service_remove` — register or deregister custom OCI services
- **MCP: `service_start` / `service_stop` support custom services** — previously only worked for built-in services
- **MCP: `.junie/guidelines.md`** — `lerd mcp:inject` now writes a lerd context section into Junie's guidelines file (merged, not overwritten) so JetBrains Junie has the same tool knowledge as Claude Code
- **Web UI: tab persistence** — active tab (Sites, Services, System) is now stored in the URL hash (`/#services`) so refreshing the browser returns to the same tab

### Fixed

- MCP skill content updated with all new tools, workflows, and architecture notes

---

## [0.5.9] — 2026-03-20

### Added

- **`lerd node:install <version>`** — install a Node.js version globally via fnm
- **`lerd node:uninstall <version>`** — uninstall a Node.js version via fnm
- **Node.js card in System tab** — lists all installed Node versions with an inline install form; replaces the install form that was previously in the Services tab
- **`lerd php:rebuild` now restarts containers** — automatically restarts all FPM containers after rebuilding images instead of printing manual instructions

### Fixed

- **`lerd tray` not opening after update** — `install.sh --update` was not copying the `lerd-tray` helper binary alongside `lerd`
- **`laravel new` and other PHP CLI tools now work end-to-end** — the PHP-FPM container image now includes Composer and Node.js/npm so subprocesses spawned by PHP (e.g. `composer create-project`, `npm install`) resolve correctly inside the container
- **`composer` and `laravel` global tools found inside container** — `lerd php` now passes the correct `HOME` and `COMPOSER_HOME` env vars and includes the Composer global bin dir in PATH so globally installed tools like the Laravel installer are found
- **Node/npm/npx shims work inside containers** — shims now use `fnm` directly (statically linked, works in Alpine) instead of calling `lerd` (glibc binary, incompatible with Alpine musl)
- **Shims use absolute paths** — `php`, `composer`, `node`, `npm`, `npx` shims now reference their binaries by absolute path, eliminating PATH-dependent failures in subprocess contexts

---

## [0.5.4] — 2026-03-19

### Added

- **Custom services**: users can now define arbitrary OCI-based services without recompiling. Config lives at `~/.config/lerd/services/<name>.yaml`.
  - `lerd service add [file.yaml]` — add from a YAML file or inline flags (`--name`, `--image`, `--port`, `--env`, `--env-var`, `--data-dir`, `--detect-key`, `--detect-prefix`, `--init-exec`, `--init-container`, `--dashboard`, `--description`)
  - `lerd service remove <name>` — stop (if running), remove quadlet and config; data directory preserved
  - `lerd service list` — shows built-in and custom services with a `[custom]` type column
  - `lerd service start/stop` — works for custom services
  - `lerd start` / `lerd stop` — includes installed custom services
  - `lerd env` — auto-detects custom services via `env_detect`, applies `env_vars`, runs `site_init.exec`
  - `lerd status` — includes custom services in the `[Services]` section
  - Web UI services tab — shows custom services with start/stop and dashboard link
  - System tray — shows custom services (slot pool expanded from 7 to 20)
- **`{{site}}` / `{{site_testing}}` placeholders** in `env_vars` and `site_init.exec` — substituted with the project site handle at `lerd env` time
- **`site_init`** YAML block — runs a `sh -c` command inside the service container once per project when `lerd env` detects the service (for DB/collection creation, user setup, etc.)
- **`dashboard`** field — shows an "Open" button in the web UI when the service is active; dashboard URLs for built-ins (Mailpit, MinIO, Meilisearch) moved from hardcoded JS to the API response
- **README simplified** — now a slim landing page pointing to the docs site
- **Docs updated** — `docs/usage/services.md` extended with full custom services reference

### Fixed

- Custom service data directory is now created automatically before starting
- `lerd service remove` now checks unit status before stopping — skips stop if not running, and aborts removal if stop fails

---

## [0.5.3] — 2026-03-19

### Fixed

- **Tray not restarting after `lerd update`**: `lerd install` was killing the tray with `pkill` but only relaunching it when `lerd-tray.service` was enabled. If the tray was started directly (`lerd tray`), it was killed and never restarted. Now tracks whether the tray was running before the kill and relaunches it directly when systemd is not managing it.

---

## [0.5.2] — 2026-03-19

### Fixed

- `lerd db:create` and `lerd db:shell` were missing from the binary — `cmd/lerd/main.go` was not staged in the v0.5.1 commit

---

## [0.5.1] — 2026-03-19

### Added

- **`lerd db:create [name]`** / **`lerd db create [name]`**: creates a database and a `<name>_testing` database in one command. Name resolution: explicit argument → `DB_DATABASE` from `.env` → project name (site registry or directory). Reports "already exists" instead of failing when a database is present. Available for both MySQL and PostgreSQL.
- **`lerd db:shell`** / **`lerd db shell`**: opens an interactive MySQL (`mysql -uroot -plerd`) or PostgreSQL (`psql -U postgres`) shell inside the service container, connecting to the project's database automatically. Replaces the need to run `podman exec --tty lerd-mysql mysql …` manually.

### Changed

- **`lerd env` now creates a `<name>_testing` database** alongside the main project database when setting up MySQL or PostgreSQL. Both databases report "already exists" if they were previously created.

---

## [0.5.0] — 2026-03-19

### Added

- **System tray applet** (`lerd tray`): a desktop tray icon for KDE, GNOME (with AppIndicator extension), waybar, and other SNI-compatible environments. The applet detaches from the terminal automatically and polls `http://127.0.0.1:7073` every 5 seconds. Menu includes:
  - 🟢/🔴 overall running status with per-component nginx and DNS indicators
  - **Open Dashboard** — opens the web UI
  - **Start / Stop Lerd** toggle
  - **Services section** — lists all active services with 🟢/🔴 status; clicking a service starts or stops it
  - **PHP section** — lists all installed PHP versions; current global default is marked ✔; clicking switches the global default via `lerd use`
  - **Autostart at login** toggle — enables or disables `lerd-autostart.service`
  - **Check for update** — polls GitHub; if a newer version is found the item changes to "⬆ Update to vX.Y.Z" and clicking opens a terminal with a confirmation prompt before running `lerd update`
  - **Stop Lerd & Quit** — runs `lerd stop` then exits the tray
- **`--mono` flag** for `lerd tray`: defaults to `true` (white monochrome icon); pass `--mono=false` for the red colour icon
- **`lerd autostart tray enable/disable`**: registers/removes `lerd-tray.service` as a user systemd unit that starts the tray on graphical login
- **`lerd start` starts the tray**: if `lerd-tray.service` is enabled it is started via systemd; otherwise, if no tray process is already running, `lerd tray` is launched directly
- **`make build-nogui`**: headless build (`CGO_ENABLED=0 -tags nogui`) for CI or servers; `lerd tray` returns a clear error instead of failing to link

### Changed

- **Build now requires CGO and `libappindicator3`** (`libappindicator-gtk3` on Arch, `libappindicator3-dev` on Debian/Ubuntu, `libappindicator-gtk3-devel` on Fedora). The `make build` target sets `CGO_ENABLED=1 -tags legacy_appindicator` automatically.
- **`lerd-autostart.service`** now declares `After=graphical-session.target` so the tray (which needs a display) is available when `lerd start` runs at login.
- **Web UI update flow**: the "Update" button has been removed. When an update is available the UI now shows `vX.Y.Z available — run lerd update in a terminal`. The `/api/update` endpoint has been removed. This avoids silent failures caused by `sudo` steps in `lerd install` that require a TTY.
- **`/api/status`** now includes a `php_default` field with the global default PHP version, used by the tray to mark the active version with ✔.

---

## [0.4.3] — 2026-03-19

### Fixed

- **DNS broken after install on Fedora (and other NM + systemd-resolved systems)**: the NetworkManager dispatcher script and `ConfigureResolver()` were calling `resolvectl domain $IFACE ~test`, which caused systemd-resolved to mark the interface as `Default Route: no`. This meant queries for anything outside `.test` (i.e. all internet DNS) had no route and were refused. Fixed by also passing `~.` as a routing domain in both places — the interface now handles `.test` specifically via lerd's dnsmasq and remains the default route for all other queries.
- **`.test` DNS fails after reboot/restart**: `lerd start` was calling `resolvectl dns` to point systemd-resolved at lerd-dns (port 5300) immediately after the container unit became active — but dnsmasq inside the container wasn't ready to accept connections yet. systemd-resolved would try port 5300, fail, mark it as a bad server, and fall back to the upstream DNS for the rest of the session. Fixed by waiting up to 10 seconds for port 5300 to accept TCP connections before calling `ConfigureResolver()`.
- **Clicking a site URL after disabling HTTPS still opened the HTTPS version**: the nginx HTTP→HTTPS redirect was a `301` (permanent), which browsers cache indefinitely. After disabling HTTPS, the browser would serve the cached redirect instead of hitting the server. Changed to `302` (temporary) so browsers always check the server, and disabling HTTPS takes effect immediately.

---

## [0.4.2] — 2026-03-19

### Changed

- **`lerd setup` detects the correct asset build command from `package.json`**: instead of always suggesting `npm run build`, the setup step now reads `scripts` from `package.json` and picks the first available candidate in priority order: `build` (Vite / default), `production` (Laravel Mix), `prod`. The step label reflects the detected command (e.g. `npm run production`). If none of the candidates exist, the build step is omitted from the selector.

---

## [0.4.1] — 2026-03-19

### Fixed

- **`lerd status` TLS certificate check**: `certExpiry` was passing raw PEM bytes directly to `x509.ParseCertificate`, which expects DER-encoded bytes. The fix decodes the PEM block first, so certificate expiry is read correctly and sites no longer show "cannot read cert" when the cert file exists and is valid.

---

## [0.4.0] — 2026-03-19

### Added

- **Xdebug toggle** (`lerd xdebug on/off [version]`): enables or disables Xdebug per PHP version by rebuilding the FPM image with Xdebug installed and configured (`mode=debug`, `start_with_request=yes`, `client_host=host.containers.internal`, port 9003). The FPM container is restarted automatically. `lerd xdebug status` shows enabled/disabled for all installed versions.
- **`lerd fetch [version...]`**: pre-builds PHP FPM images for the specified versions (or all supported: 8.1–8.5) so the first `lerd use <version>` is instant. Skips versions whose images already exist.
- **`lerd db:import <file.sql>`** / **`lerd db:export [-o file]`**: import or export a SQL dump using the project's `.env` DB settings. Supports MySQL/MariaDB (`lerd-mysql`) and PostgreSQL (`lerd-postgres`). Also available as `lerd db import` / `lerd db export`.
- **`lerd share [site]`**: exposes the current site publicly via ngrok or Expose. Auto-detects which tunnel tool is installed; use `--ngrok` or `--expose` to force one. Forwards to the local nginx port with the correct `Host` header so nginx routes to the right vhost.
- **`lerd setup`**: interactive project bootstrap command — presents a checkbox list of steps (composer install, npm ci, lerd env, lerd mcp:inject, php artisan migrate, php artisan db:seed, npm run build, lerd secure, lerd open) with smart defaults based on project state. `lerd link` always runs first (mandatory, not in the list) to ensure the site is registered with the correct PHP version before any subsequent step. `--all` / `-a` runs everything without prompting (CI-friendly); `--skip-open` skips opening the browser.

### Fixed

- **PHP version detection order**: `composer.json` `require.php` now takes priority over `.php-version`, so projects declaring `"php": "^8.4"` in `composer.json` automatically use PHP 8.4 even if a stale `.php-version` file says otherwise. Explicit `.lerd.yaml` overrides still take top priority.
- **`lerd link` preserves HTTPS**: re-linking a site that was already secured now regenerates the SSL vhost (not an HTTP vhost), so `https://` continues to work after a re-link.
- **`lerd link` preserves `secured` flag**: re-linking no longer resets a secured site to `secured: false`.
- **`lerd secure` / `lerd unsecure` directory name resolution**: sites in directories with real TLDs (e.g. `astrolov.com`) are now resolved correctly by path lookup, so the commands no longer error with "site not found" when the directory name differs from the registered site name.

---

## [0.3.0] — 2026-03-18

### Added

- `lerd env` command: copies `.env.example` → `.env` if missing, detects which services the project uses, applies lerd connection values, starts required services, generates `APP_KEY` if missing, and sets `APP_URL` to the registered `.test` domain
- `lerd unsecure [name]` command: removes the mkcert TLS cert and reverts the site to HTTP
- `lerd secure` and `lerd unsecure` now automatically update `APP_URL` in the project's `.env` to `https://` or `http://` respectively
- `lerd install` now installs a `/etc/sudoers.d/lerd` rule granting passwordless `resolvectl dns/domain/revert` — required for the autostart service which cannot prompt for a sudo password
- PHP FPM images now include the `gmp` extension
- **MCP server** (`lerd mcp`): JSON-RPC 2.0 stdio server exposing lerd as a Model Context Protocol tool provider for AI assistants (Claude Code, JetBrains Junie, and any MCP-compatible client). Tools: `artisan`, `sites`, `service_start`, `service_stop`, `queue_start`, `queue_stop`, `logs`
- **`lerd mcp:inject`**: writes `.mcp.json`, `.claude/skills/lerd/SKILL.md`, and `.junie/mcp/mcp.json` into a project directory. Merges into existing `mcpServers` configs — other servers (e.g. `laravel-boost`, `herd`) are preserved unchanged
- **UI: queue worker toggle** in the Sites tab — amber toggle to start/stop the queue worker per site; spinner while toggling; error text on failure; **logs** link opens the live log drawer for that worker when running
- **UI: Unlink button** in the Sites tab — small red-bordered button that confirms, calls `POST /api/sites/{domain}/unlink`, and removes the site from the table client-side immediately
- **`lerd unlink` parked-site behaviour**: unlinking a site under a parked directory now marks it as `ignored` in the registry instead of removing it, preventing the watcher from re-registering it on next scan. Running `lerd link` in the same directory clears the flag. Non-parked sites are still removed from the registry entirely
- `GET /api/sites` filters out ignored sites so they are invisible in the UI
- `queue:start` and `queue:stop` are now also available as API actions via `POST /api/sites/{domain}/queue:start` and `POST /api/sites/{domain}/queue:stop`, enabling UI and MCP control

### Fixed

- DNS `.test` routing now works correctly after autostart: `resolvectl revert` is called before re-applying per-interface DNS settings so systemd-resolved resets the current server to `127.0.0.1:5300`; previously, resolved would mark lerd-dns as failed during boot (before it started) then fall back to the upstream DNS for all queries including `.test`, causing NXDOMAIN on every `.test` lookup
- `fnm install` no longer prints noise to the terminal when a Node version is already installed

### Changed

- `lerd start` and `lerd stop` now start/stop containers in parallel — startup is noticeably faster on multi-container setups
- `lerd start` now re-applies DNS resolver config on every invocation, ensuring `.test` routing is always correct after reboot or network changes
- `lerd park` now skips already-registered sites instead of overwriting them, preserving settings such as TLS status and custom PHP version
- `lerd install` completion message now shows both `http://lerd.test` and `http://127.0.0.1:7073` as fallback
- Composer is now stored as `composer.phar`; the `composer` shim runs it via `lerd php`
- Autostart service now declares `After=network-online.target` and runs at elevated priority (`Nice=-10`)

---

## [0.2.0] — 2026-03-17

### Changed

- UI completely redesigned: dark theme inspired by Laravel.com with near-black background, red accents, and top navbar replacing the sidebar
- Light / Auto / Dark theme toggle added to the navbar; preference persists in localStorage

---

## [0.1.0] — 2026-03-17

Initial release.

### Added

**Core**
- Single static Go binary built with Cobra
- XDG-compliant config (`~/.config/lerd/`) and data (`~/.local/share/lerd/`) directories
- Global config at `~/.config/lerd/config.yaml` with sensible defaults
- Per-project `.lerd.yaml` override support
- Linux distro detection (Arch, Debian/Ubuntu, Fedora, openSUSE)
- Build metadata injected at compile time: version, commit SHA, build date

**Site management**
- `lerd park [dir]` — auto-discover and register all Laravel projects in a directory
- `lerd link [name]` — register the current directory as a named site
- `lerd unlink` — remove a site and clean up its vhost
- `lerd sites` — tabular view of all registered sites

**PHP**
- `lerd install` — one-time setup: directories, Podman network, binary downloads, DNS, nginx
- `lerd use <version>` — set the global PHP version
- `lerd isolate <version>` — pin PHP version per-project via `.php-version`
- `lerd php:list` — list installed static PHP binaries
- PHP version resolution order: `.php-version` → `.lerd.yaml` → `composer.json` → global default

**Node**
- `lerd isolate:node <version>` — pin Node version per-project via `.node-version`
- Node version resolution order: `.nvmrc` → `.node-version` → `package.json engines.node` → global default
- fnm bundled for Node version management

**TLS**
- `lerd secure [name]` — issue a locally-trusted mkcert certificate for a site
- Automatic HTTPS vhost generation
- mkcert CA installed into system trust store on `lerd install`

**Services**
- `lerd service start|stop|restart|status|list` — manage optional services
- Bundled services: MySQL 8.0, Redis 7, PostgreSQL 16, Meilisearch v1.7, MinIO

**Infrastructure**
- All containers run rootless on a dedicated `lerd` Podman network
- Nginx and PHP-FPM as Podman Quadlet containers (auto-managed by systemd)
- dnsmasq container for `.test` TLD resolution via NetworkManager
- fsnotify-based watcher daemon (`lerd-watcher.service`) for auto-discovery of new projects

**Diagnostics**
- `lerd status` — health overview: DNS, nginx, PHP-FPM containers, services, cert expiry
- `lerd dns:check` — verify `.test` resolution

**Lifecycle**
- `lerd update` — self-update from latest GitHub release (atomic binary swap)
- `lerd uninstall` — stop all containers, remove units, binary, PATH entry, optionally data
- Shell completion via `lerd completion bash|zsh|fish`

---

[0.6.0]: https://github.com/geodro/lerd/compare/v0.5.16...v0.6.0
[0.5.16]: https://github.com/geodro/lerd/compare/v0.5.15...v0.5.16
[0.5.15]: https://github.com/geodro/lerd/compare/v0.5.14...v0.5.15
[0.5.14]: https://github.com/geodro/lerd/compare/v0.5.13...v0.5.14
[0.5.13]: https://github.com/geodro/lerd/compare/v0.5.12...v0.5.13
[0.5.12]: https://github.com/geodro/lerd/compare/v0.5.11...v0.5.12
[0.5.11]: https://github.com/geodro/lerd/compare/v0.5.9...v0.5.11
[0.5.9]: https://github.com/geodro/lerd/compare/v0.5.4...v0.5.9
[0.5.4]: https://github.com/geodro/lerd/compare/v0.5.3...v0.5.4
[0.5.3]: https://github.com/geodro/lerd/compare/v0.5.2...v0.5.3
[0.5.2]: https://github.com/geodro/lerd/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/geodro/lerd/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/geodro/lerd/compare/v0.4.3...v0.5.0
[0.4.3]: https://github.com/geodro/lerd/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/geodro/lerd/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/geodro/lerd/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/geodro/lerd/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/geodro/lerd/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/geodro/lerd/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/geodro/lerd/releases/tag/v0.1.0
