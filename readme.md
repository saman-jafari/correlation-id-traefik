# Correlation ID plugin for Traefik

A [Traefik](https://traefik.io/) middleware that guarantees every request carries a correlation ID header — reusing the client's if present, or generating a [UUID v7](https://uuid7.com/) if not — so a request can be traced across all the services it touches.

[![Build](https://github.com/saman-jafari/correlation-id-traefik/actions/workflows/go-cross.yml/badge.svg)](https://github.com/saman-jafari/correlation-id-traefik/actions/workflows/go-cross.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/saman-jafari/correlation-id-traefik)](https://goreportcard.com/report/github.com/saman-jafari/correlation-id-traefik)
[![Latest release](https://img.shields.io/github/v/tag/saman-jafari/correlation-id-traefik?label=release)](https://github.com/saman-jafari/correlation-id-traefik/tags)
[![License](https://img.shields.io/github/license/saman-jafari/correlation-id-traefik)](LICENSE)

## TL;DR

- **What:** adds a correlation ID header to every incoming request.
- **Behavior:** header already set by the client → kept as-is; header missing → new UUID v7 generated.
- **Default header:** `X-Correlation-Id` (override with `headerName`).
- **Compatible with:** Traefik v2 and v3.

```yaml
# traefik static config
experimental:
  plugins:
    correlation:
      moduleName: github.com/saman-jafari/correlation-id-traefik
      version: v1.0.1
```
```yaml
# attach to a service (docker labels)
labels:
  - traefik.http.middlewares.correlation.plugin.correlation.headerName=X-Correlation-Id
```

## How it works

For each request the middleware reads the configured header:

1. **Present** → the value is forwarded unchanged, keeping the ID stable across hops.
2. **Absent** → a UUID v7 is generated and set on the header before the request continues.

UUID v7 is time-ordered (its leading bits are a Unix millisecond timestamp), so IDs sort by creation time — convenient for logs and databases.

## Installation

### As a published plugin (production)

Register the plugin in Traefik's **static** configuration, pinned to a released tag:

```yaml
# traefik.yml
experimental:
  plugins:
    correlation:
      moduleName: github.com/saman-jafari/correlation-id-traefik
      version: v1.0.1
```

Or via CLI flags:

```yaml
command:
  - "--experimental.plugins.correlation.moduleName=github.com/saman-jafari/correlation-id-traefik"
  - "--experimental.plugins.correlation.version=v1.0.1"
```

### As a local plugin (development)

Mount the source into Traefik and reference it as a local plugin — see [`docker-compose.yml`](docker-compose.yml) for a complete, runnable example.

```yaml
command:
  - "--experimental.localPlugins.correlation.moduleName=github.com/saman-jafari/correlation-id-traefik"
volumes:
  - "./:/plugins-local/src/github.com/saman-jafari/correlation-id-traefik/"
```

## Configuration

| Option       | Type   | Default            | Description                                          |
|--------------|--------|--------------------|------------------------------------------------------|
| `headerName` | string | `X-Correlation-Id` | Header read from the request and set on it if empty. |

Header names are case-insensitive, so `X-Correlation-Id`, `x-correlation-id`, and `X-CORRELATION-ID` are equivalent.

### Apply to a single service

```yaml
# docker labels
labels:
  - traefik.http.middlewares.correlation.plugin.correlation.headerName=X-Correlation-Id
  - traefik.http.routers.whoami.middlewares=correlation
```

### Apply to every route on an entrypoint

```yaml
command:
  - "--entrypoints.web.http.middlewares=correlation@docker"
  - "--entrypoints.websecure.http.middlewares=correlation@docker"
```

## Try it locally

```shell
docker compose up -d

# header echoed back unchanged
curl -s http://whoami.localhost/ -H 'X-Correlation-Id: my-trace-id' | grep -i correlation

# a fresh UUID v7 is generated
curl -s http://whoami.localhost/ | grep -i correlation
```

## Development

No Go toolchain? Everything runs in Docker (see commands below).

```shell
make            # lint + test + yaegi_test (run before every commit)
make test       # go test -v -cover ./...
make yaegi_test # run tests through the Yaegi interpreter
make lint       # golangci-lint (v2 config)
make vendor     # re-vendor dependencies after editing go.mod
```

> [!IMPORTANT]
> Traefik executes plugins with the [Yaegi](https://github.com/traefik/yaegi) interpreter, **not** the Go compiler. `go test` passing does not prove the plugin loads in Traefik — always run `make yaegi_test`.

Run the toolchain without installing Go locally:

```shell
# tests + Yaegi
docker run --rm -v "$PWD":/app -w /app golang:1.23 sh -c 'go test -cover ./...'

# lint
docker run --rm -v "$PWD":/app -w /app golangci/golangci-lint:v2.12.2 golangci-lint run
```

Pinned toolchain versions (Go, golangci-lint, Yaegi) live in [`.github/workflows/main.yml`](.github/workflows/main.yml).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, checks, and the release process. Changes are tracked in [CHANGELOG.md](CHANGELOG.md).

## Compatibility

| Component | Version        |
|-----------|----------------|
| Traefik   | v2.x, v3.x     |
| Go        | 1.23+          |

## License

[Apache-2.0](LICENSE)
