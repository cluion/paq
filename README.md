# paq

A unified CLI for querying installed packages across multiple package managers.

[繁體中文](README.zh-TW.md)

## Why?

Developers typically have multiple package managers installed, each with different commands:

```bash
brew list          # Homebrew
npm list -g        # npm
snap list          # Snap
flatpak list       # Flatpak
```

**paq** gives you one consistent interface:

```bash
paq brew
paq npm
paq snap
paq flatpak
```

## Features

- **Single binary** — cross-platform (macOS, Linux, Windows)
- **Zero config** — auto-detects available package managers
- **Two output formats** — table (default) and JSON (`--json`)
- **Extensible** — add new providers without modifying core code

## Install

### Homebrew

```bash
brew install cluion/tap/paq
```

### Go install

```bash
go install github.com/cluion/paq@latest
```

### From source

```bash
git clone https://github.com/cluion/paq.git
cd paq
make build
sudo make install   # installs to /usr/local/bin/paq
```

If you prefer not to use `sudo`, add an alias to your shell config (`~/.zshrc` or `~/.bashrc`):

```bash
alias paq="/path/to/paq/bin/paq"
```

### Download binary

Download the latest release from [Releases](https://github.com/cluion/paq/releases).

## Update

**Go:**

```bash
go install github.com/cluion/paq@latest
```

**Homebrew:**

```bash
brew upgrade cluion/tap/paq
```

**From source:**

```bash
cd paq
git pull
make build
sudo make install   # if alias, just rebuild — it takes effect immediately
```

## Releasing

Maintainers can trigger a release by pushing a tag. GitHub Actions will automatically build, create a Release, and update the Homebrew formula:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Prerequisite: create the [cluion/homebrew-tap](https://github.com/cluion/homebrew-tap) repository on GitHub first.

## Usage

### List detected package sources

```bash
paq list
```

```
Source  Display Name  Status
brew    Homebrew      available
npm     npm           available
snap    Snap          not found
flatpak Flatpak       not found
Total: 4
```

### Query a specific source

```bash
paq brew
```

```
Name        Version   Description
git         2.40.0    Distributed version control system
node        20.1.0    Platform built on V8
curl        8.1.0
Total: 3
```

### JSON output

```bash
paq brew --json
```

```json
{
  "source": "brew",
  "count": 3,
  "packages": [
    { "name": "git", "version": "2.40.0", "desc": "Distributed version control system" },
    { "name": "node", "version": "20.1.0", "desc": "Platform built on V8" },
    { "name": "curl", "version": "8.1.0" }
  ]
}
```

### Version

```bash
paq version
```

```
paq dev (commit: abc1234, built: 2026-04-06)
```

## Supported Package Managers

| Provider | Command | Platforms |
|----------|---------|-----------|
| Homebrew | `paq brew` | macOS |
| npm | `paq npm` | macOS, Linux, Windows |
| Snap | `paq snap` | Linux |
| Flatpak | `paq flatpak` | Linux |

## Development

### Prerequisites

- Go 1.22+
- golangci-lint

### Build

```bash
make build
```

### Test

```bash
make test
```

### Lint

```bash
make lint
```

## License

[MIT](LICENSE)
