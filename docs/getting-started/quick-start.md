# Quick Start

Get a Laravel project running locally in three commands.

```bash
# 1. Park your projects directory
#    Any Laravel project inside is auto-registered
lerd park ~/Lerd

# 2. Visit your project in a browser
#    ~/Lerd/my-app  →  http://my-app.test

# 3. Check everything is running
lerd status
```
{ .annotate }

1. `lerd park` registers the directory with the watcher service. Every subdirectory that looks like a Laravel project gets a `.test` domain automatically.
2. No `/etc/hosts` edits needed — DNS is handled by dnsmasq running in a Podman container.
3. `lerd status` shows a health summary: DNS, nginx, PHP-FPM containers, services, and cert expiry.

That's it. Nginx is serving your project through PHP-FPM, all inside Podman containers on the `lerd` network.

---

## First project bootstrap

For a freshly cloned project, use `lerd setup` to run all the standard setup steps in one go:

```bash
cd ~/Lerd/my-app
lerd setup
```

A checkbox list appears with all available steps pre-selected. Toggle steps with space, confirm with enter, and watch them run sequentially.

See [Project Setup](../features/project-setup.md) for the full details.

---

## Web UI

The dashboard is available at **`http://127.0.0.1:7073`** once Lerd is installed. It gives you a visual overview of all your sites, services, and system health.

See [Web UI](../features/web-ui.md) for details.
