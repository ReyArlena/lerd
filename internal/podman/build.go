package podman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// BuildFPMImage builds the lerd PHP-FPM image for the given version if it doesn't exist.
// Prints build output to stdout so the user can see progress.
func BuildFPMImage(version string) error {
	short := strings.ReplaceAll(version, ".", "")
	imageName := "lerd-php" + short + "-fpm:local"

	// Skip if image already exists
	checkCmd := exec.Command("podman", "image", "exists", imageName)
	if checkCmd.Run() == nil {
		return nil
	}

	fmt.Printf("\n  Building PHP %s image (first time, may take a few minutes)...\n", version)

	containerfileTmpl, err := GetQuadletTemplate("lerd-php-fpm.Containerfile")
	if err != nil {
		return err
	}
	containerfile := strings.ReplaceAll(containerfileTmpl, "{{.Version}}", version)

	tmp, err := os.MkdirTemp("", "lerd-php-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	cfPath := tmp + "/Containerfile"
	if err := os.WriteFile(cfPath, []byte(containerfile), 0644); err != nil {
		return err
	}

	cmd := exec.Command("podman", "build", "-t", imageName, "-f", cfPath, tmp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("building PHP %s image: %w", version, err)
	}

	fmt.Printf("  PHP %s image built successfully.\n", version)
	return nil
}
