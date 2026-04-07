package repository

import (
	"fmt"
	"path/filepath"
	"sort"

	"curriculum/internal/store"

	"golang.org/x/mod/semver"
)

// SkillInfo describes a skill version available in the central repository.
type SkillInfo struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
	Latest   string   `json:"latest"`
}

// Push copies the skill directory at srcDir into the central repository under
// ~/.curriculum/repository/<name>/<version>/.
func Push(name, version, srcDir string) (string, error) {
	repoDir, err := store.RepositoryDir()
	if err != nil {
		return "", err
	}
	dest := filepath.Join(repoDir, name, version)
	if err := store.EnsureDir(dest); err != nil {
		return "", fmt.Errorf("create repository dir: %w", err)
	}
	if err := store.CopyDir(srcDir, dest); err != nil {
		return "", fmt.Errorf("copy skill: %w", err)
	}
	return dest, nil
}

// Install copies the skill from the central repository to destBase/<name>/.
// If version is empty, the latest available version is used.
func Install(name, version, destBase string) (string, error) {
	ver, err := resolveVersion(name, version)
	if err != nil {
		return "", err
	}

	repoDir, err := store.RepositoryDir()
	if err != nil {
		return "", err
	}
	src := filepath.Join(repoDir, name, ver)
	if !store.Exists(src) {
		return "", fmt.Errorf("skill %q version %s not found in repository", name, ver)
	}

	dest := filepath.Join(destBase, name)
	if err := store.EnsureDir(dest); err != nil {
		return "", fmt.Errorf("create install dir: %w", err)
	}
	if err := store.CopyDir(src, dest); err != nil {
		return "", fmt.Errorf("copy skill: %w", err)
	}
	return dest, nil
}

// Remove deletes the skill directory at skillsBase/<name>/.
// Returns os.ErrNotExist if the skill is not installed there.
func Remove(name, skillsBase string) error {
	target := filepath.Join(skillsBase, name)
	if !store.Exists(target) {
		return fmt.Errorf("skill %q not found at %s", name, target)
	}
	return store.DeleteDir(target)
}

// List returns all skills available in the central repository.
func List() ([]*SkillInfo, error) {
	repoDir, err := store.RepositoryDir()
	if err != nil {
		return nil, err
	}
	names, err := store.ListDir(repoDir)
	if err != nil {
		return nil, err
	}
	var result []*SkillInfo
	for _, name := range names {
		versions, err := store.ListDir(filepath.Join(repoDir, name))
		if err != nil {
			return nil, err
		}
		if len(versions) == 0 {
			continue
		}
		latest := latestVersion(versions)
		result = append(result, &SkillInfo{
			Name:     name,
			Versions: versions,
			Latest:   latest,
		})
	}
	return result, nil
}

// resolveVersion returns version if non-empty, otherwise the latest version for name.
func resolveVersion(name, version string) (string, error) {
	if version != "" {
		return version, nil
	}
	repoDir, err := store.RepositoryDir()
	if err != nil {
		return "", err
	}
	versions, err := store.ListDir(filepath.Join(repoDir, name))
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("skill %q not found in repository", name)
	}
	return latestVersion(versions), nil
}

// latestVersion returns the highest semver from a slice of version strings.
// Falls back to the last element (lexicographic) if none parse as semver.
func latestVersion(versions []string) string {
	// semver package requires a "v" prefix.
	sv := make([]string, 0, len(versions))
	for _, v := range versions {
		candidate := "v" + v
		if semver.IsValid(candidate) {
			sv = append(sv, v)
		}
	}
	if len(sv) == 0 {
		// Fall back: return last lexicographic entry.
		sorted := make([]string, len(versions))
		copy(sorted, versions)
		sort.Strings(sorted)
		return sorted[len(sorted)-1]
	}
	sort.Slice(sv, func(i, j int) bool {
		return semver.Compare("v"+sv[i], "v"+sv[j]) < 0
	})
	return sv[len(sv)-1]
}
