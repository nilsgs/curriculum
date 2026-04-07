package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"curriculum/internal/manifest"
	"curriculum/internal/repository"
	"curriculum/internal/store"

	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an installed skill",
		Long: `Remove a skill from .agents/skills/<name>.

Use --global to remove from ~/.agents/skills/ instead.
Use --no-save to skip removing the dependency from .curriculum.`,
		Args: cobra.ExactArgs(1),
		RunE: runRemove,
	}
	cmd.Flags().Bool("global", false, "remove from ~/.agents/skills/ instead of .agents/skills/")
	cmd.Flags().Bool("no-save", false, "do not remove the dependency from .curriculum")
	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	name := args[0]
	global, _ := cmd.Flags().GetBool("global")
	noSave, _ := cmd.Flags().GetBool("no-save")
	save := !noSave

	var skillsBase string
	if global {
		base, err := store.GlobalSkillsDir()
		if err != nil {
			return err
		}
		skillsBase = base
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		skillsBase = filepath.Join(cwd, ".agents", "skills")
	}

	if err := repository.Remove(name, skillsBase); err != nil {
		return notFoundErr("skill", name)
	}

	if save {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not update %s: %v\n", manifest.FileName, err)
		} else {
			m, repoRoot, err := manifest.Load(cwd)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not load %s: %v\n", manifest.FileName, err)
			} else {
				m.RemoveDep(name)
				if err := manifest.Save(repoRoot, m); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not save %s: %v\n", manifest.FileName, err)
				}
			}
		}
	}

	output(cmd, map[string]string{"name": name, "status": "removed"}, func(w io.Writer) {
		fmt.Fprintf(w, "Removed %s\n", name)
	})
	return nil
}
