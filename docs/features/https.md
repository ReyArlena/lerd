# HTTPS / TLS

Lerd uses [mkcert](https://github.com/FiloSottile/mkcert) — a locally-trusted CA that your browser will accept without warnings.

```bash
cd ~/Lerd/my-app
lerd secure
# Issues a cert for my-app.test, regenerates the SSL vhost, reloads nginx
# Updates APP_URL=https://my-app.test in .env if it exists
# Visit https://my-app.test — no certificate warning

lerd unsecure
# Removes the cert, switches back to HTTP vhost
# Updates APP_URL=http://my-app.test in .env if it exists
```

Certificates are stored in `~/.local/share/lerd/certs/sites/`.

---

## From the Web UI

The Sites tab has an HTTPS toggle per site — clicking it runs `lerd secure` or `lerd unsecure` inline and updates the vhost without touching the terminal.

---

## How it works

1. `lerd install` generates a local CA with mkcert and installs it into the system trust store (NSS databases for Chrome/Firefox, and the system root store).
2. `lerd secure <site>` issues a certificate signed by that CA for the site's `.test` domain.
3. The nginx vhost is regenerated to listen on port 443 with the new cert, and port 80 redirects to HTTPS (302, not 301, so the redirect is not cached by browsers).
4. `APP_URL` in the project's `.env` is updated to `https://`.
