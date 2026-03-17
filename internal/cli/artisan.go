package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/geodro/lerd/internal/config"
	phpDet "github.com/geodro/lerd/internal/php"
	"github.com/spf13/cobra"
)

// NewArtisanCmd returns the artisan command — runs php artisan in the project's container.
func NewArtisanCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "artisan [args...]",
		Short:              "Run php artisan in the project's container",
		Example:            "  lerd artisan migrate\n  lerd artisan tinker\n  lerd artisan make:controller Foo",
		DisableFlagParsing: true,
		RunE:               runArtisan,
	}
}

func runArtisan(_ *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	version, err := phpDet.DetectVersion(cwd)
	if err != nil {
		cfg, cfgErr := config.LoadGlobal()
		if cfgErr != nil {
			return fmt.Errorf("cannot detect PHP version: %w", err)
		}
		version = cfg.PHP.DefaultVersion
	}

	short := strings.ReplaceAll(version, ".", "")
	container := "lerd-php" + short + "-fpm"

	cmdArgs := []string{"exec", "-it", "-w", cwd, container, "php", "artisan"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("podman", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
