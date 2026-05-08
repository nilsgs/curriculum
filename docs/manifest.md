# curriculum Manifest

The `.curriculum` file is a JSON manifest that declares skills a repository
provides and skills it consumes.

## Example

```json
{
  "version": "1.0.0",
  "skills": [
    {
      "name": "my-skill"
    },
    {
      "name": "custom-skill",
      "path": "agent-tools/custom-skill"
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

## Fields

`version` is required. It must be a semantic version string and is applied to
all provided skills when `cur push` publishes them.

`skills` is optional. It lists skills this repository provides.

`dependencies` is optional. It lists skills this repository consumes.

## Provided Skills

Each `skills` entry supports:

- `name`: required skill name. It must match the `SKILL.md` frontmatter name.
- `path`: optional path from the repository root. Defaults to `skills/<name>`.

Example default path:

```json
{
  "skills": [
    { "name": "my-skill" }
  ]
}
```

The skill is read from:

```text
skills/my-skill/
```

Example custom path:

```json
{
  "skills": [
    {
      "name": "my-skill",
      "path": "agent-tools/my-skill"
    }
  ]
}
```

The skill is read from:

```text
agent-tools/my-skill/
```

## Dependencies

Each `dependencies` entry supports:

- `name`: required skill name in the central repository.
- `version`: optional semantic version. If omitted, install resolves the latest
  available version at install time.

Pinned dependency:

```json
{
  "dependencies": [
    {
      "name": "my-skill",
      "version": "1.0.0"
    }
  ]
}
```

Floating dependency:

```json
{
  "dependencies": [
    {
      "name": "my-skill"
    }
  ]
}
```

## Skill Directory

A skill directory must contain `SKILL.md`:

```text
my-skill/
  SKILL.md
  scripts/
  references/
  assets/
```

Only `SKILL.md` is required. Other files and folders are copied as part of the
skill.

`SKILL.md` must include frontmatter:

```yaml
---
name: my-skill
description: "What this skill teaches agents and when to use it"
---

# My Skill

Instructions for agents.
```

The frontmatter `name` must match the manifest entry name.
