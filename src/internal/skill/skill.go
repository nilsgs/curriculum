package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const skillFile = "SKILL.md"

// Frontmatter holds the parsed YAML frontmatter from a SKILL.md file.
type Frontmatter struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	License     string            `yaml:"license,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

var validName = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// ValidateName checks that a skill name follows the agentskills.io naming rules.
func ValidateName(name string) error {
	if len(name) == 0 || len(name) > 64 {
		return fmt.Errorf("skill name must be 1-64 characters")
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("skill name %q must be lowercase alphanumeric with hyphens (no leading/trailing/consecutive hyphens)", name)
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("skill name %q must not contain consecutive hyphens", name)
	}
	return nil
}

// ValidateDir checks that dir is a valid skill directory:
// - SKILL.md must exist
// - SKILL.md frontmatter name must match expectedName
// - directory base name must match expectedName
func ValidateDir(dir, expectedName string) error {
	skillPath := filepath.Join(dir, skillFile)
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return fmt.Errorf("skill directory %q is missing %s", dir, skillFile)
	}

	fm, err := ParseFrontmatter(skillPath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", skillPath, err)
	}

	if fm.Name != expectedName {
		return fmt.Errorf("%s frontmatter name %q does not match expected name %q", skillFile, fm.Name, expectedName)
	}

	if filepath.Base(dir) != expectedName {
		return fmt.Errorf("skill directory base name %q does not match skill name %q", filepath.Base(dir), expectedName)
	}

	return nil
}

// ParseFrontmatter extracts and parses the YAML frontmatter from a SKILL.md file.
func ParseFrontmatter(path string) (*Frontmatter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Expect opening ---
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, fmt.Errorf("missing YAML frontmatter opening ---")
	}

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(strings.Join(lines, "\n")), &fm); err != nil {
		return nil, fmt.Errorf("parse YAML frontmatter: %w", err)
	}

	if fm.Name == "" {
		return nil, fmt.Errorf("frontmatter missing required 'name' field")
	}
	if fm.Description == "" {
		return nil, fmt.Errorf("frontmatter missing required 'description' field")
	}

	return &fm, nil
}
