package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"curriculum/internal/manifest"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a .curriculum manifest in the current directory",
		Long:  "Creates a .curriculum file with empty skills and dependencies. Fails if one already exists.",
		RunE:  runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	if manifest.ExistsIn(cwd) {
		return fmt.Errorf("already initialized in this directory (found %s)", manifest.FileName)
	}

	m := manifest.Empty()
	if err := manifest.Save(cwd, m); err != nil {
		return fmt.Errorf("write %s: %w", manifest.FileName, err)
	}

	path := filepath.Join(cwd, manifest.FileName)
	output(cmd, map[string]string{"path": path, "version": m.Version}, func(w io.Writer) {
		fmt.Fprintf(w, "Initialized %s\n", path)
	})
	return nil
}
