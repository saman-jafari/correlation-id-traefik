# Contributing

Thanks for your interest in improving this Traefik plugin. This guide covers the local setup, the checks that must pass, and how a release reaches the Traefik Plugin Catalog.

## Prerequisites

- Go **1.23+** — or Docker only (every command below has a Docker equivalent, no local Go needed).
- `make`.
- Optional: [golangci-lint](https://golangci-lint.run/) v2.12.2 and [Yaegi](https://github.com/traefik/yaegi) v0.16.1 if you prefer running them natively. Pinned versions live in [`.github/workflows/main.yml`](.github/workflows/main.yml).

## Development workflow

```shell
make            # lint + test + yaegi_test — run this before every commit
make test       # go test -v -cover ./...
make yaegi_test # run tests through the Yaegi interpreter
make lint       # golangci-lint (v2 config)
make vendor     # re-vendor after editing go.mod (see "Dependencies")
```

Without a local Go toolchain:

```shell
docker run --rm -v "$PWD":/app -w /app golang:1.23 sh -c 'go test -cover ./...'
docker run --rm -v "$PWD":/app -w /app golangci/golangci-lint:v2.12.2 golangci-lint run
```

> [!IMPORTANT]
> Traefik executes plugins with the **Yaegi interpreter**, not the Go compiler. A passing `go test` does **not** prove the plugin loads in Traefik. `make yaegi_test` is mandatory — CI runs it, and the Plugin Catalog re-runs the `.traefik.yml` `testData` under Yaegi at ingest time.

## Dependencies

Dependencies are **vendored and committed** (`vendor/`) — this is a hard requirement for Traefik plugins. After changing an import or `go.mod`:

```shell
go mod tidy
go mod vendor
```

CI fails if `go.mod`, `go.sum`, or `vendor/` drift from a clean `tidy` + `vendor`.

## Coding standards

- `golangci-lint run` must report **0 issues** (config: [`.golangci.yml`](.golangci.yml), v2 schema).
- Keep test coverage — the package currently sits at 100%.
- Do not rename the plugin entry points required by Traefik's contract: `CreateConfig`, `New`, and the `ServeHTTP` method.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/) — the history already follows it (`feat:`, `fix:`). This keeps the changelog and version bumps predictable.

## Pull requests

1. Branch from `master`.
2. Make the change; add or update tests.
3. Run `make` (lint + test + yaegi_test) — all green.
4. Update [`CHANGELOG.md`](CHANGELOG.md) under `## [Unreleased]` and any affected docs.
5. Open the PR against `master`; CI must pass.

## Releasing

Versions follow [SemVer](https://semver.org/). The Traefik Plugin Catalog resolves versions from **git tags via the Go module proxy** — the tag *is* the published version. There is no submission form and no version field to edit; publishing a tag is the entire release.

1. Move `## [Unreleased]` entries in `CHANGELOG.md` under a new `## [X.Y.Z]` heading and update the compare links.
2. Commit, then tag and push:
   ```shell
   git tag vX.Y.Z
   git push origin master --tags
   ```
3. The catalog polls GitHub **once a day** and picks up the new tag automatically (up to ~24h). If ingest fails, the catalog **opens an issue on this repo** and pauses until it is closed.

Catalog listing requirements — all already satisfied, keep them intact:

- repository is public and not a fork;
- the `traefik-plugin` GitHub topic is set;
- [`.traefik.yml`](.traefik.yml) exists with valid `testData`;
- a valid `go.mod` at the repo root;
- dependencies are vendored and committed;
- the release is a semver git tag.
