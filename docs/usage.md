# curriculum Usage

curriculum manages reusable agent skills between repositories. The `cur` CLI
uses a `.curriculum` manifest in each repository and a local central store under
`~/.curriculum` by default.

## Initialize A Repository

```sh
cur init
```

This creates a `.curriculum` manifest with an initial version, no provided
skills, and no dependencies.

## Publish A Skill

Create a skill directory containing `SKILL.md`, then declare it in
`.curriculum`:

```json
{
  "version": "1.0.0",
  "skills": [
    { "name": "my-skill" }
  ],
  "dependencies": []
}
```

Publish all declared skills:

```sh
cur push
```

Publish one declared skill:

```sh
cur push my-skill
```

Published skills are copied to:

```text
~/.curriculum/repository/<skill-name>/<version>/
```

## Install Skills

Install all dependencies declared in `.curriculum`:

```sh
cur install
```

Install the latest version of one skill and save it to `.curriculum`:

```sh
cur install my-skill
```

Install a specific version:

```sh
cur install my-skill@1.0.0
```

Install without updating `.curriculum`:

```sh
cur install my-skill --no-save
```

Install to personal skills instead of repo-local skills:

```sh
cur install my-skill --global
```

Default install target:

```text
.agents/skills/<skill-name>/
```

Global install target:

```text
~/.agents/skills/<skill-name>/
```

## Remove Skills

```sh
cur remove my-skill
cur remove my-skill --no-save
cur remove my-skill --global
```

Repo-local removal deletes `.agents/skills/<skill-name>/` and updates
`.curriculum` unless `--no-save` is used.

## List Available Skills

```sh
cur list
cur list --json
```

`cur list` reads from the local central store and shows available skill names
and versions.

## JSON Output

Commands accept `--json` for machine-readable output:

```sh
cur list --json
cur push --json
```

Errors are written to stderr. JSON output is written to stdout.

## Storage

Default central store:

```text
~/.curriculum/
  repository/
    <skill-name>/
      <version>/
        SKILL.md
        scripts/
        references/
        assets/
```

Set `CURRICULUM_HOME` to use another storage root.
