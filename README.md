# `curriculum` — Manage Agentic Skills Between Repositories

`curriculum` is a CLI tool for publishing and consuming agentic skills following the [agentskills.io specification](https://agentskills.io/specification). A skill is a reusable agent instruction (SKILL.md + optional scripts/assets) that can be shared across repositories.

**Key features:**
- **Push skills** to a central repository (`~/.curriculum/repository/`)
- **Install dependencies** into any repo (`cur install`)
- **Semver versioning** — each skill is versioned independently
- **Declarative manifest** — `.curriculum` JSON declares what your repo provides and consumes
- **Multi-target installs** — install to repo-local (`.agents/skills/`) or personal (`~/.agents/skills/`)

---

## Installation

### From source (macOS, Linux, Windows)

```bash
git clone https://github.com/nskut/curriculum.git
cd curriculum

# Linux / macOS
./install.sh

# Windows PowerShell
.\install.ps1
```

Installs the `cur` binary to your `$GOPATH/bin` (add to `$PATH` if needed).

### Verify installation

```bash
cur --version
cur --help
```

---

## Quick Start

### 1. Initialize a repo

```bash
cd my-repo
cur init
```

Creates `.curriculum` with empty `skills` and `dependencies`.

### 2. Declare a skill your repo provides

Create a skill directory and SKILL.md following [agentskills.io spec](https://agentskills.io/specification):

```bash
mkdir -p skills/my-skill
cat > skills/my-skill/SKILL.md << 'EOF'
---
name: my-skill
description: "Teaches agents how to do X"
---

# My Skill

Detailed instructions...
EOF
```

Update `.curriculum`:

```json
{
  "version": "1.0.0",
  "skills": [
    { "name": "my-skill" }
  ],
  "dependencies": []
}
```

### 3. Push to central repository

```bash
cur push
```

Pushes all skills to `~/.curriculum/repository/my-skill/1.0.0/`.

### 4. Use in another repo

```bash
cd another-repo
cur init
cur install my-skill  # Installs and adds to .curriculum dependencies
cur install my-skill --no-save  # Install without updating .curriculum
```

---

## Command Reference

### Global Flags

| Flag | Description |
|---|---|
| `--json` | Output machine-readable JSON to stdout |
| `--version` | Print version and exit |
| `-h, --help` | Print help for the command |

**Exit codes:**
- `0` — success
- `1` — error
- `2` — entity not found

---

### `cur init`

Initialize `.curriculum` in the current directory.

```bash
cur init
```

Creates `.curriculum` with:
```json
{
  "version": "0.1.0",
  "skills": [],
  "dependencies": []
}
```

**Options:**
- None

---

### `cur push [<name>]`

Push skills to the central repository at `~/.curriculum/repository/`.

```bash
# Push all declared skills
cur push

# Push a specific skill
cur push my-skill

# JSON output
cur push --json
```

Validates each skill:
- `SKILL.md` must exist
- Frontmatter `name` must match manifest entry `name`
- Directory basename must match skill name

Each skill is pushed to:
```
~/.curriculum/repository/<name>/<version>/
```

where `<version>` comes from the `.curriculum` top-level `version` field.

---

### `cur install [<name> [@<version>]]`

Install skills from the central repository into `.agents/skills/`.

```bash
# Install all dependencies from .curriculum
cur install

# Install a single skill (defaults to latest version)
cur install my-skill

# Install a specific version
cur install my-skill@1.0.0

# Install to personal skills (~/.agents/skills/ instead of .agents/skills/)
cur install my-skill --global

# Install without updating .curriculum
cur install my-skill --no-save
```

**Behavior:**
- No args: reads `dependencies[]` from `.curriculum`, installs all
- With `<name>`: installs from central repo
- With `@version`: installs that specific semver (if omitted, installs latest)
- `--global`: installs to `~/.agents/skills/` instead of `.agents/skills/`
- `--no-save`: skips updating `.curriculum` `dependencies[]`

---

### `cur remove <name>`

Remove an installed skill.

```bash
# Remove from repo-local .agents/skills/ and update .curriculum
cur remove my-skill

# Remove from personal ~/.agents/skills/
cur remove my-skill --global

# Remove without updating .curriculum
cur remove my-skill --no-save
```

**Behavior:**
- Deletes the skill directory and removes the entry from `.curriculum` `dependencies[]`
- `--global`: removes from `~/.agents/skills/` instead of `.agents/skills/`
- `--no-save`: skips updating `.curriculum` `dependencies[]`

---

### `cur list`

List all skills available in the central repository.

```bash
# Human-readable table
cur list

# JSON output
cur list --json
```

Shows skill names and available versions. Latest version is marked with `*`.

---

## Configuration

### `.curriculum` Manifest

The `.curriculum` file is a JSON manifest that declares:
1. **Skills this repo provides** — source code lives in `.skills/<name>/` by convention
2. **Skills this repo consumes** — installed into `.agents/skills/`

#### Schema

```json
{
  "version": "1.0.0",
  "skills": [
    {
      "name": "my-skill",
      "path": "skills/my-skill"
    },
    {
      "name": "another",
      "path": "custom/location/another"
    }
  ],
  "dependencies": [
    {
      "name": "external-skill",
      "version": "2.1.0"
    },
    {
      "name": "latest-skill"
    }
  ]
}
```

**`version`** (required)
- Semver string (e.g., `"1.0.0"`)
- Applied to *all* skills when pushing
- Bump this to publish a new version of all skills

**`skills[]`** (optional)
- Array of skills this repo provides
- `name` — skill identifier; must match SKILL.md frontmatter `name`
- `path` — relative path from repo root; defaults to `skills/<name>` if omitted

**`dependencies[]`** (optional)
- Array of skills this repo consumes
- `name` — skill identifier in central repository
- `version` — semver to install; optional, defaults to latest if omitted

#### Path Resolution

**Skills you provide:**
- Default location: `skills/<name>/`
- Override per entry with `path`
- Example: `"path": "agent-tools/my-skill"` → reads from `agent-tools/my-skill/`

**Skills you consume:**
- Always installed to `.agents/skills/<name>/` (repo-local)
- Or `~/.agents/skills/<name>/` with `--global` flag

---

## Skill Directory Structure

Skills follow [agentskills.io spec](https://agentskills.io/specification):

```
my-skill/
├── SKILL.md          # Required: YAML frontmatter + instructions
├── scripts/          # Optional: executable code (bash, Python, etc.)
├── references/       # Optional: additional docs, reference files
├── assets/           # Optional: templates, images, data files
└── ...               # Any additional files or directories
```

### SKILL.md Format

Must contain YAML frontmatter:

```yaml
---
name: my-skill
description: "What this skill teaches agents and when to use it"
license: "MIT"
compatibility: "Requires Node.js 16+"
metadata:
  category: "dev-tools"
  tags: "testing, automation"
---

# My Skill

Instructions and detailed guidance for agents...
```

**Required fields:**
- `name` — max 64 chars, lowercase alphanumeric + hyphens only
- `description` — max 1024 chars, non-empty

**Optional fields:**
- `license`, `compatibility`, `metadata`

---

## Environment Variables

### `CURRICULUM_HOME`

Override the default central repository location (`~/.curriculum`).

```bash
export CURRICULUM_HOME=/custom/path/to/curriculum
cur push              # Uses /custom/path/to/curriculum/repository/
cur list              # Lists from /custom/path/to/curriculum/repository/
```

If not set, defaults to `~/.curriculum`.

---

## Examples

### Example 1: Publish a skill

```bash
mkdir -p my-repo/skills/python-testing
cat > my-repo/skills/python-testing/SKILL.md << 'EOF'
---
name: python-testing
description: "Teaches agents pytest patterns, fixtures, and best practices"
---

# Python Testing

When to use this skill:
- Writing unit tests with pytest
- Testing async code
- Parameterized tests
...
EOF

cd my-repo
cat > .curriculum << 'EOF'
{
  "version": "1.0.0",
  "skills": [
    { "name": "python-testing" }
  ],
  "dependencies": []
}
EOF

cur push
# → Pushed to ~/.curriculum/repository/python-testing/1.0.0/
```

### Example 2: Install a skill

```bash
cd another-project
cur init

# Install the latest version — also saved to .curriculum automatically
cur install python-testing
# → Installed to .agents/skills/python-testing/
# → Added to .curriculum dependencies

# Verify it's there
ls -la .agents/skills/python-testing/
```

### Example 3: Personal (global) skills

```bash
# Install a skill to your personal skills directory
cur install python-testing --global
# → Installed to ~/.agents/skills/python-testing/
# → Available to all repos on this machine
```

### Example 4: Version pinning

```bash
# Install a specific version
cur install python-testing@1.0.0
# → .curriculum dependency: { "name": "python-testing", "version": "1.0.0" }

# Install latest of a different skill
cur install another-skill
# → .curriculum dependency: { "name": "another-skill" }
# → Resolves to highest available version at install time
```

---

## Central Repository Layout

The central repository is organized by skill name and version:

```
~/.curriculum/
└── repository/
    ├── python-testing/
    │   ├── 1.0.0/
    │   │   ├── SKILL.md
    │   │   ├── scripts/
    │   │   └── references/
    │   └── 1.1.0/
    │       ├── SKILL.md
    │       └── ...
    ├── my-skill/
    │   ├── 0.5.0/
    │   └── 1.0.0/
    └── ...
```

---

## Troubleshooting

### "no .curriculum found (run 'cur init' first)"

Run `cur init` in the repository root to create a `.curriculum` manifest.

### "skill 'X' not declared in .curriculum skills"

When pushing, the skill must be listed in `.curriculum` `skills[]`.

### "skill directory must match SKILL.md frontmatter name"

The directory basename must match the frontmatter `name` field:
- Directory: `skills/my-skill/`
- `SKILL.md`: `name: my-skill` ✓

### "skill not found in repository"

The skill hasn't been pushed yet, or it's in a different central repo. Check:
```bash
cur list | grep skill-name
```

### Installing to the wrong location

- Default: `.agents/skills/<name>/` (repo-local)
- Personal: `~/.agents/skills/<name>/` (with `--global`)

---

## Contributing

Prerequisites for development:

- Go 1.26+
- Task v3, installed from the official instructions: <https://taskfile.dev/docs/installation>
- Docker or Podman for `task smoke` and `task ci`

```bash
task test     # native Go tests
task build    # build local binary into dist/
task install  # copy dist/cur to ~/.curriculum/bin
task smoke    # run Smoko specs; .smokorc builds the test image
task ci       # run test, build, and smoke
task cross    # build the full OS/architecture matrix into dist/
task clean    # remove dist/
```

Equivalent native Go test command for debugging Task itself:

```bash
cd src && go test ./... -v -count=1
```

See `AGENTS.md` for the development workflow.

---
## License

See LICENSE file.
