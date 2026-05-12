![Banner](img/banner.png)

# curriculum

curriculum is a CLI for publishing and installing agent skills between
repositories.

The `cur` command pushes skills from a repository into a local central store,
installs skills into other repositories, and records dependencies in a
`.curriculum` manifest.

## Install

Prerequisites:

- Go 1.26+
- Git

Linux / macOS:

```sh
git clone https://github.com/nilsgs/curriculum.git
cd curriculum
./install.sh
```

Windows PowerShell:

```powershell
git clone https://github.com/nilsgs/curriculum.git
cd curriculum
.\install.ps1
```

The installer builds from source, copies `cur` to `~/.curriculum/bin`, and
updates your user `PATH` where supported.

## Quick Start

```sh
cur init
cur push
cur install my-skill
cur list
```

## Usage

Use `--help` for the full command surface:

```sh
cur --help
cur install --help
cur push --help
```

Common workflows:

```sh
cur init
cur push my-skill
cur install my-skill
cur install my-skill --global
cur list --json
```

Skills are read from the paths declared in `.curriculum`. Installed repo-local
skills go to `.agents/skills/`; personal skills installed with `--global` go to
`~/.agents/skills/`.

By default, curriculum stores published skills under `~/.curriculum`. Set
`CURRICULUM_HOME` to use another location.

## Docs

- [Expanded usage](docs/usage.md)
- [Manifest reference](docs/manifest.md)
- [Agent skill guide](skills/curriculum/SKILL.md)

## Development

Prerequisites:

- Go 1.26+
- Task v3: <https://taskfile.dev/docs/installation>
- Docker or Podman for `task smoke` and `task ci`

Common tasks:

```sh
task test     # run native Go tests
task build    # build the local binary into dist/
task install  # build and copy the binary to the user install directory
task smoke    # run Smoko specs
task ci       # run test, build, and smoke
task cross    # build the full OS/architecture matrix into dist/
task clean    # remove dist/
```

Smoke specs use tags for focused runs, for example `smoko run specs/ --tag json`.

The version is read from `VERSION` and stamped into the binary at build time.

## License

MIT. See [LICENSE](LICENSE).
