# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.0] - 2026-04-09

### Added

- Table styling with rounded borders and cell padding
- `Type` field in package output (e.g. `formula` / `cask` for Homebrew)
- Color-coded status in `paq list` (green for available, gray for not found)
- Dynamic column visibility — Type and Description columns only appear when relevant

### Changed

- Improved `paq list` table with full rounded border style
- Improved `paq brew/npm/...` table with lightweight header separator
- Long package names and versions are now truncated with `…`
- Removed Description column from package listing tables

## [1.3.0] - 2026-04-08

### Added

- APT/dpkg provider (`paq apt`) for Debian/Ubuntu systems
- Test suite for apt provider

### Changed

- README Go version requirement updated from 1.22+ to 1.26+ to match go.mod

## [1.2.1] - 2026-04-07

### Fixed

- `paq upgrade` now compares versions for all methods (Homebrew, go install, self-update)
- Fix errcheck lint error in version command
- CI Go version updated to 1.26 to match go.mod

### Added

- CodeQL security analysis in CI
- govulncheck vulnerability scanning in CI
- Coverage threshold check (80%) in CI

## [1.2.0] - 2026-04-07

### Fixed

- `paq upgrade` now compares versions before self-update download
- Write permission check before self-update to prevent silent failures

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
