package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"curriculum/internal/store"
)

const FileName = ".curriculum"

// Manifest is the per-repo curriculum configuration.
type Manifest struct {
	Version      string       `json:"version"`
	Skills       []SkillEntry `json:"skills"`
	Dependencies []DepEntry   `json:"dependencies"`
}

// SkillEntry describes a skill this repo provides.
type SkillEntry struct {
	Name string `json:"name"`
	// Path is the relative path from the repo root to the skill directory.
	// Defaults to skills/<name> if empty.
	Path string `json:"path,omitempty"`
}

// ResolvePath returns the effective path for the skill entry.
// If Path is set it is returned as-is; otherwise the default skills/<Name> is used.
func (e SkillEntry) ResolvePath() string {
	if e.Path != "" {
		return e.Path
	}
	return filepath.Join("skills", e.Name)
}

// DepEntry describes a skill this repo consumes.
type DepEntry struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// Empty returns a new Manifest with empty slices (suitable for cur init).
func Empty() *Manifest {
	return &Manifest{
		Version:      "0.1.0",
		Skills:       []SkillEntry{},
		Dependencies: []DepEntry{},
	}
}

// Save writes the manifest to the given directory.
func Save(dir string, m *Manifest) error {
	return store.SaveJSON(filepath.Join(dir, FileName), m)
}

// Load walks up from startDir looking for a .curriculum file.
// Returns the manifest and the directory where it was found.
func Load(startDir string) (*Manifest, string, error) {
	dir := startDir
	for {
		path := filepath.Join(dir, FileName)
		if store.Exists(path) {
			var m Manifest
			if err := store.LoadJSON(path, &m); err != nil {
				return nil, "", fmt.Errorf("read %s: %w", path, err)
			}
			if m.Skills == nil {
				m.Skills = []SkillEntry{}
			}
			if m.Dependencies == nil {
				m.Dependencies = []DepEntry{}
			}
			return &m, dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, "", fmt.Errorf("no %s found (run 'cur init' first)", FileName)
}

// ExistsIn returns true if a .curriculum file exists in the given directory.
func ExistsIn(dir string) bool {
	return store.Exists(filepath.Join(dir, FileName))
}

// FindSkill returns the SkillEntry with the given name, or an error if not found.
func (m *Manifest) FindSkill(name string) (SkillEntry, error) {
	for _, s := range m.Skills {
		if s.Name == name {
			return s, nil
		}
	}
	return SkillEntry{}, fmt.Errorf("skill %q not declared in %s skills", name, FileName)
}

// FindDep returns the DepEntry with the given name, or an error if not found.
func (m *Manifest) FindDep(name string) (DepEntry, error) {
	for _, d := range m.Dependencies {
		if d.Name == name {
			return d, nil
		}
	}
	return DepEntry{}, fmt.Errorf("dependency %q not declared in %s", name, FileName)
}

// UpsertDep adds or updates a dependency entry in the manifest.
func (m *Manifest) UpsertDep(name, version string) {
	for i, d := range m.Dependencies {
		if d.Name == name {
			m.Dependencies[i].Version = version
			return
		}
	}
	m.Dependencies = append(m.Dependencies, DepEntry{Name: name, Version: version})
}

// RemoveDep removes a dependency entry by name. Returns true if it was present.
func (m *Manifest) RemoveDep(name string) bool {
	for i, d := range m.Dependencies {
		if d.Name == name {
			m.Dependencies = append(m.Dependencies[:i], m.Dependencies[i+1:]...)
			return true
		}
	}
	return false
}

// ResolveStartDir returns the current working directory, panicking on error.
func ResolveStartDir() (string, error) {
	return os.Getwd()
}
