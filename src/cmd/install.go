package cmd

import (
	"fmt"
	"io"
	"os"

	"curriculum/internal/manifest"
	"curriculum/internal/repository"
	"curriculum/internal/store"

	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [<name> [@<version>]]",
		Short: "Install skills from the central repository",
		Long: `Install skills from ~/.curriculum/repository/ into .agents/skills/.

Without arguments, installs all dependencies declared in .curriculum.
With a name (and optional @version), installs that single skill.

Use --global to install into ~/.agents/skills/ instead.
Use --save to record the skill as a dependency in .curriculum.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runInstall,
	}
	cmd.Flags().Bool("global", false, "install into ~/.agents/skills/ instead of .agents/skills/")
	cmd.Flags().Bool("save", false, "add or update the dependency in .curriculum")
	return cmd
}

type installResult struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Destination string `json:"destination"`
}

func runInstall(cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")
	save, _ := cmd.Flags().GetBool("save")

	destBase, err := resolveInstallBase(global)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return installAllDeps(cmd, destBase, global, save)
	}

	nameArg, version := splitNameVersion(args[0])
	dest, err := repository.Install(nameArg, version, destBase)
	if err != nil {
		return err
	}

	// Resolve the version that was actually installed.
	resolvedVersion := version
	if resolvedVersion == "" {
		infos, _ := repository.List()
		for _, info := range infos {
			if info.Name == nameArg {
				resolvedVersion = info.Latest
				break
			}
		}
	}

	if save {
		if err := saveDepToManifest(nameArg, resolvedVersion); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not update %s: %v\n", manifest.FileName, err)
		}
	}

	result := installResult{Name: nameArg, Version: resolvedVersion, Destination: dest}
	output(cmd, result, func(w io.Writer) {
		fmt.Fprintf(w, "Installed %s@%s → %s\n", result.Name, result.Version, result.Destination)
	})
	return nil
}

func installAllDeps(cmd *cobra.Command, destBase string, global, save bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	m, repoRoot, err := manifest.Load(cwd)
	if err != nil {
		return err
	}

	if len(m.Dependencies) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No dependencies declared in .curriculum")
		return nil
	}

	var results []installResult
	for _, dep := range m.Dependencies {
		dest, err := repository.Install(dep.Name, dep.Version, destBase)
		if err != nil {
			return fmt.Errorf("install %q: %w", dep.Name, err)
		}
		resolvedVersion := dep.Version
		if resolvedVersion == "" {
			infos, _ := repository.List()
			for _, info := range infos {
				if info.Name == dep.Name {
					resolvedVersion = info.Latest
					break
				}
			}
		}
		results = append(results, installResult{Name: dep.Name, Version: resolvedVersion, Destination: dest})
	}

	if save {
		for _, r := range results {
			m.UpsertDep(r.Name, r.Version)
		}
		if err := manifest.Save(repoRoot, m); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not update %s: %v\n", manifest.FileName, err)
		}
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

// resolveInstallBase returns the destination base directory for installs.
func resolveInstallBase(global bool) (string, error) {
	if global {
		return store.GlobalSkillsDir()
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	return fmt.Sprintf("%s/.agents/skills", cwd), nil
}

// splitNameVersion splits "name@version" into name and version.
// If there is no "@", version is returned as empty string.
func splitNameVersion(arg string) (name, version string) {
	for i, c := range arg {
		if c == '@' {
			return arg[:i], arg[i+1:]
		}
	}
	return arg, ""
}

// saveDepToManifest adds or updates a dependency in the nearest .curriculum.
func saveDepToManifest(name, version string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, repoRoot, err := manifest.Load(cwd)
	if err != nil {
		return err
	}
	m.UpsertDep(name, version)
	return manifest.Save(repoRoot, m)
}
