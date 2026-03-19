# PHP

## Commands

| Command | Description |
|---|---|
| `lerd use <version>` | Set the global PHP version and build the FPM image if needed |
| `lerd isolate <version>` | Pin PHP version for cwd — writes `.php-version` |
| `lerd php:list` | List all installed PHP-FPM versions |
| `lerd php:rebuild` | Force-rebuild all installed PHP-FPM images (run after `lerd update` if needed) |
| `lerd fetch [version...]` | Pre-build PHP FPM images for the given (or all supported) versions so first use isn't slow |
| `lerd php [args...]` | Run PHP in the project's container |
| `lerd artisan [args...]` | Run `php artisan` in the project's container |
| `lerd xdebug on [version]` | Enable Xdebug for a PHP version — rebuilds the FPM image and restarts the container |
| `lerd xdebug off [version]` | Disable Xdebug — rebuilds without Xdebug and restarts |
| `lerd xdebug status` | Show Xdebug enabled/disabled for all installed PHP versions |

If no version is given, the version is resolved from the current directory (`.php-version` or `composer.json`, falling back to the global default).

---

## Version resolution

When serving a request, Lerd picks the PHP version for a project in this order:

1. `.lerd.yaml` in the project root — `php_version` field (explicit lerd override)
2. `composer.json` — `require.php` constraint (e.g. `^8.4` → `8.4`)
3. `.php-version` file in the project root (plain text, e.g. `8.2`)
4. Global default in `~/.config/lerd/config.yaml`

To pin a project permanently:

```bash
cd ~/Lerd/my-app
lerd isolate 8.2
# writes .php-version: 8.2 — commit this if you like
```

To change the global default:

```bash
lerd use 8.4
```

---

## Xdebug

??? info "Xdebug configuration values"
    Xdebug is configured with:

    - `xdebug.mode=debug`
    - `xdebug.start_with_request=yes`
    - `xdebug.client_host=host.containers.internal` (reaches your host IDE from the container)
    - `xdebug.client_port=9003`

    Set your IDE to listen on port `9003`. In VS Code, the default PHP Debug configuration works without changes. In PhpStorm, set **Settings → PHP → Debug → Debug port** to `9003`.
