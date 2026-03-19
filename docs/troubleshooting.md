# Troubleshooting

??? bug "`.test` domains not resolving"
    Run the DNS check first:

    ```bash
    lerd dns:check
    ```

    If it fails, restart NetworkManager and check again:

    ```bash
    sudo systemctl restart NetworkManager
    lerd dns:check
    ```

    On systems using systemd-resolved (Ubuntu), check that the per-interface DNS configuration was applied:

    ```bash
    resolvectl status
    # Look for your default interface — it should show 127.0.0.1:5300 as DNS server
    # and ~test as a routing domain
    ```

??? bug "Nginx not serving a site"
    Check that nginx and the PHP-FPM container are running, then inspect the generated vhost:

    ```bash
    lerd status                         # check nginx and FPM are running
    podman logs lerd-nginx              # nginx error log
    cat ~/.local/share/lerd/nginx/conf.d/my-app.test.conf   # check generated vhost
    ```

??? bug "PHP-FPM container not running"
    Check the systemd unit status and logs:

    ```bash
    systemctl --user status lerd-php84-fpm
    systemctl --user start lerd-php84-fpm
    podman logs lerd-php84-fpm
    ```

    If the image is missing (e.g. after `podman rmi`):

    ```bash
    lerd php:rebuild
    ```

??? bug "Permission denied on port 80/443"
    Rootless Podman cannot bind to ports below 1024 by default. Allow it:

    ```bash
    sudo sysctl -w net.ipv4.ip_unprivileged_port_start=80
    # Make permanent:
    echo 'net.ipv4.ip_unprivileged_port_start=80' | sudo tee /etc/sysctl.d/99-lerd.conf
    ```

    `lerd install` sets this automatically, but it may need to be re-applied after a kernel update.

??? bug "Watcher service not running"
    The watcher auto-discovers new projects in parked directories. If sites aren't being picked up:

    ```bash
    systemctl --user status lerd-watcher
    systemctl --user start lerd-watcher
    ```

??? bug "HTTPS certificate warning in browser"
    The mkcert CA must be installed in your browser's trust store. Ensure `certutil` / `nss-tools` is installed, then re-run `lerd install`:

    - Arch: `sudo pacman -S nss`
    - Debian/Ubuntu: `sudo apt install libnss3-tools`
    - Fedora: `sudo dnf install nss-tools`

    After installing the package, run `lerd install` again to register the CA.
