package dns

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/geodro/lerd/internal/config"
)

const dnsmasqConf = `# Lerd DNS configuration
address=/.test/127.0.0.1
port=53
`

const nmDnsConf = `[main]
dns=dnsmasq
`

const nmDnsmasqConf = `server=/test/127.0.0.1#5300
`

const resolvedDropin = `[Resolve]
DNS=127.0.0.1:5300
Domains=~test
`

// isFileContent returns true if the file at path already contains exactly content.
func isFileContent(path string, content []byte) bool {
	existing, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return string(existing) == string(content)
}

// isSystemdResolvedActive returns true if systemd-resolved is the active DNS resolver.
func isSystemdResolvedActive() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "systemd-resolved")
	if err := cmd.Run(); err != nil {
		return false
	}
	// Also check that /etc/resolv.conf points to the stub resolver
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "127.0.0.53") || strings.Contains(string(data), "systemd-resolved")
}

// Setup writes DNS configuration for .test resolution and restarts the resolver.
// On systemd-resolved systems (Ubuntu etc.) it uses a resolved drop-in.
// On NetworkManager-only systems it uses NM's embedded dnsmasq.
func Setup() error {
	if err := WriteDnsmasqConfig(config.DnsmasqDir()); err != nil {
		return fmt.Errorf("writing lerd dnsmasq config: %w", err)
	}

	if isSystemdResolvedActive() {
		return setupSystemdResolved()
	}
	return setupNetworkManager()
}

// setupSystemdResolved configures systemd-resolved to forward .test to port 5300.
func setupSystemdResolved() error {
	dropin := "/etc/systemd/resolved.conf.d/lerd.conf"

	if isFileContent(dropin, []byte(resolvedDropin)) {
		return nil
	}

	fmt.Println("  [sudo required] Configuring systemd-resolved for .test DNS resolution")

	if err := sudoWriteFile(dropin, []byte(resolvedDropin)); err != nil {
		return fmt.Errorf("writing resolved drop-in: %w", err)
	}

	cmd := exec.Command("sudo", "systemctl", "restart", "systemd-resolved")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restarting systemd-resolved: %w", err)
	}
	return nil
}

// setupNetworkManager configures NetworkManager's embedded dnsmasq.
func setupNetworkManager() error {
	nmConfFile := "/etc/NetworkManager/conf.d/lerd.conf"
	nmDnsmasqFile := "/etc/NetworkManager/dnsmasq.d/lerd.conf"

	if isFileContent(nmConfFile, []byte(nmDnsConf)) && isFileContent(nmDnsmasqFile, []byte(nmDnsmasqConf)) {
		return nil
	}

	fmt.Println("  [sudo required] Configuring NetworkManager for .test DNS resolution")

	if err := sudoWriteFile(nmConfFile, []byte(nmDnsConf)); err != nil {
		return fmt.Errorf("writing NetworkManager conf: %w", err)
	}

	if err := sudoWriteFile(nmDnsmasqFile, []byte(nmDnsmasqConf)); err != nil {
		return fmt.Errorf("writing NetworkManager dnsmasq conf: %w", err)
	}

	cmd := exec.Command("sudo", "systemctl", "restart", "NetworkManager")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restarting NetworkManager: %w", err)
	}
	return nil
}

// WriteDnsmasqConfig writes the lerd dnsmasq config to the given directory.
func WriteDnsmasqConfig(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "lerd.conf"), []byte(dnsmasqConf), 0644)
}

// sudoWriteFile writes content to a system path by writing to a temp file
// then using sudo cp, so sudo can prompt for a password on the terminal.
func sudoWriteFile(path string, content []byte) error {
	tmp, err := os.CreateTemp("", "lerd-sudo-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	dir := filepath.Dir(path)
	mkdirCmd := exec.Command("sudo", "mkdir", "-p", dir)
	mkdirCmd.Stdin = os.Stdin
	mkdirCmd.Stdout = os.Stdout
	mkdirCmd.Stderr = os.Stderr
	if err := mkdirCmd.Run(); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	cpCmd := exec.Command("sudo", "cp", tmp.Name(), path)
	cpCmd.Stdin = os.Stdin
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("cp to %s: %w", path, err)
	}
	return nil
}
