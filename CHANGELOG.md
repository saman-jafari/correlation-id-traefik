# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0]

### Changed
- Generate a UUID v7 only when the header is absent (previously generated one on every request and then discarded it when a value was already present). External behavior is unchanged.
- Raised the minimum Go version to 1.23 (`go.mod`).
- Migrated `.golangci.yml` to the golangci-lint v2 schema.
- Modernized CI: `actions/checkout@v4`, `actions/setup-go@v5`, `golangci-lint-action@v7`, pinned Go 1.23 / golangci-lint v2.12.2 / Yaegi v0.16.1, and added a module tidiness/vendoring drift check.
- `docker-compose.yml` now uses `traefik:v3.7` and drops the obsolete Compose `version` field.
- Rewrote the README (TL;DR, configuration table, development section) and expanded developer docs.

### Added
- `make yaegi_test` target and a CI step that runs the tests through the Yaegi interpreter — the same way Traefik loads the plugin at runtime.
- Tests covering ID generation, custom header names, and preservation of an incoming ID.

### Fixed
- `New` no longer mutates the caller-supplied `Config`.

## [1.0.1]

- Naming fixes.

## [1.0.0]

- Initial release: Traefik middleware that sets a correlation ID header (UUID v7) on each request, preserving an existing value when present.

[Unreleased]: https://github.com/saman-jafari/correlation-id-traefik/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/saman-jafari/correlation-id-traefik/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/saman-jafari/correlation-id-traefik/releases/tag/v1.0.0
