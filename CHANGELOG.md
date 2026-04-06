# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-04-07

### Added

- `paq upgrade` command with auto-detection of install method (Homebrew, go install, local binary, direct download)
- `paq upgrade --check` flag to check for new version without upgrading
- Self-update via direct binary download from GitHub Releases (`minio/selfupdate`)
- `make uninstall` Makefile target
- Uninstall instructions in README (English and Chinese)
- `debug.ReadBuildInfo()` fallback for version display when ldflags are not set

### Changed

- `go install` path corrected from `github.com/cluion/paq@latest` to `github.com/cluion/paq/cmd/paq@latest`
- Version display now uses `debug.ReadBuildInfo()` as fallback, fixing version info for `go install` builds

## [1.0.0] - 2026-04-06

### Added

- Initial release of paq
- Unified CLI for querying installed packages across package managers
- Homebrew provider (`paq brew`)
- npm provider (`paq npm`)
- Snap provider (`paq snap`)
- Flatpak provider (`paq flatpak`)
- `paq list` command to show detected package sources
- `paq version` command with build info injection
- Table output format (default)
- JSON output format (`--json` flag)
- Auto-detection of available package managers
- Cross-platform support (macOS, Linux, Windows)
- Test coverage >= 80%
- Makefile with build, test, lint, clean, install targets
- GoReleaser configuration for automated releases
- GitHub Actions CI/CD pipeline
