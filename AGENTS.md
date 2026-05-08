# curriculum Agent Notes

## Layout

- `src/main.go`: CLI entry point.
- `src/cmd/`: Cobra commands for `init`, `push`, `install`, `remove`, and `list`.
- `src/internal/manifest/`: `.curriculum` parsing and validation.
- `src/internal/store/`: central repository path and storage behavior.
- `skills/curriculum/`: agent-facing curriculum usage guidance.
- `specs/`: Smoko smoke specs.

## Build And Test

Use Task targets for normal validation:

```sh
task test
task build
task smoke
task ci
```

Raw Go fallback for focused debugging:

```sh
cd src
go test ./... -v -count=1
```

## Smoke Tests

Smoke specs live under `specs/` and are run with Smoko:

```sh
task smoke
```

`.smokorc` owns the image build. Do not duplicate Docker build commands in Task
or docs unless the project workflow changes.

## Contracts

- `.curriculum` is the manifest for provided skills and dependencies.
- Repo-local installs go to `.agents/skills/`.
- Global installs go to `~/.agents/skills/`.
- The central store defaults to `~/.curriculum`.
- `CURRICULUM_HOME` overrides the central store root.
- `--json` output must remain parseable on stdout.

## Documentation

- Keep `README.md` as the concise user and developer front door.
- Put expanded user workflows in `docs/usage.md`.
- Put manifest details in `docs/manifest.md`.
- Keep agent-specific skill guidance in `skills/curriculum/SKILL.md`.
- Keep Markdown docs ASCII-only unless a non-ASCII character is deliberate.

## Versioning

The version comes from `VERSION` and is stamped into builds by the Task build
scripts.
