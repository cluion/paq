# Contributing to paq

Thank you for your interest in contributing to paq! This document provides guidelines for contributions.

[繁體中文](CONTRIBUTING.zh-TW.md)

## How to Contribute

### Report Bugs

Open an [issue](https://github.com/cluion/paq/issues/new) with:

- OS and architecture
- paq version (`paq version` output)
- Steps to reproduce
- Expected vs actual behavior

### Suggest Features

Open an [issue](https://github.com/cluion/paq/issues/new) with the label `enhancement`.

### Submit Code

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Ensure all checks pass:
   ```bash
   make lint
   make test
   ```
5. Commit with [conventional commits](https://www.conventionalcommits.org/):
   ```
   feat: add cargo provider
   fix: handle empty npm output
   docs: update README
   ```
6. Push and open a Pull Request

## Development Setup

### Prerequisites

- Go 1.22+
- [golangci-lint](https://golangci-lint.run/usage/install/)

### Build & Test

```bash
make build    # compile to bin/paq
make test     # run tests with race detector and coverage
make lint     # run golangci-lint
make clean    # remove build artifacts
```

### Project Structure

```
cmd/paq/main.go              # Entry point
internal/
  cli/                        # CLI commands (root, list, query, version)
  provider/                   # Provider interface + implementations
  output/                     # Formatters (table, JSON)
Makefile
.goreleaser.yml
```

### Adding a New Provider

1. Create `internal/provider/<name>.go`
2. Implement the `Provider` interface:
   ```go
   type Provider interface {
       Name() string
       DisplayName() string
       Detect() bool
       List() ([]Package, error)
   }
   ```
3. Register in `init()`:
   ```go
   func init() {
       Register(&MyProvider{runner: defaultRunner})
   }
   ```
4. Add tests with a fake `CommandRunner`
5. Update README and CHANGELOG

### Code Style

- Follow `gofmt` and `goimports`
- Error messages should be lowercase, no trailing punctuation
- Interfaces should be small (1-3 methods)
- Keep functions under 50 lines, files under 800 lines

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
