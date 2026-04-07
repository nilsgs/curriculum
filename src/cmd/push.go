package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"curriculum/internal/manifest"
	"curriculum/internal/repository"
	"curriculum/internal/skill"

	"github.com/spf13/cobra"
)

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push [<name>]",
		Short: "Push skills to the central repository",
		Long: `Push skills declared in .curriculum to ~/.curriculum/repository/.

Without arguments, pushes all skills. With a name, pushes that skill only.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runPush,
	}
}

type pushResult struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Destination string `json:"destination"`
}

func runPush(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	m, repoRoot, err := manifest.Load(cwd)
	if err != nil {
		return err
	}

	if m.Version == "" {
		return fmt.Errorf(".curriculum has no version set")
	}

	var entries []manifest.SkillEntry
	if len(args) == 1 {
		e, err := m.FindSkill(args[0])
		if err != nil {
			return notFoundErr("skill", args[0])
		}
		entries = []manifest.SkillEntry{e}
	} else {
		if len(m.Skills) == 0 {
			return fmt.Errorf("no skills declared in %s", manifest.FileName)
		}
		entries = m.Skills
	}

	var results []pushResult
	for _, e := range entries {
		srcDir := filepath.Join(repoRoot, filepath.FromSlash(e.ResolvePath()))
		if err := skill.ValidateDir(srcDir, e.Name); err != nil {
			return fmt.Errorf("skill %q: %w", e.Name, err)
		}
		dest, err := repository.Push(e.Name, m.Version, srcDir)
		if err != nil {
			return fmt.Errorf("push skill %q: %w", e.Name, err)
		}
		results = append(results, pushResult{
			Name:        e.Name,
			Version:     m.Version,
			Destination: dest,
		})
	}

	output(cmd, results, func(w io.Writer) {
		rows := make([][]string, len(results))
		for i, r := range results {
			rows[i] = []string{r.Name, r.Version, r.Destination}
		}
		printTable(w, []string{"SKILL", "VERSION", "DESTINATION"}, rows)
	})
	return nil
}
