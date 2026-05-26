# ConfigGen — AGENTS.md

## Quick start

```powershell
go build -o configgen.exe ./cmd/configgen
.\configgen.exe -addr=:8080
```

Verify: `curl http://localhost:8080/health` → `{"status":"ok"}`

## Commands

| Action | Command |
|--------|---------|
| Build | `go build -o configgen.exe ./cmd/configgen` |
| Test all | `go test ./...` |
| Test one package | `go test ./internal/infrastructure/generators/` |
| Static analysis | `go vet ./...` |
| Coverage | `go test -cover ./...` |

No formatter or linter config beyond `go fmt` / `go vet`.

## Architecture

Clean Architecture, four layers:

```
cmd/configgen/main.go          — DI assembly: store → generator → usecase → server
internal/domain/               — entities (ConfigRequest, ConfigResult, ConfigRecord) + Store interface
internal/usecase/              — business logic: applyDefaults, GenerateAndSave, GetRecord, ListRecords
internal/infrastructure/       — concrete implementations (generators, SQLite store)
internal/presentation/         — Gin HTTP handlers + request parsing
```

Dependency direction: `presentation → usecase → domain ← infrastructure`

## Templates are compiled into the binary

All templates live under `internal/infrastructure/generators/templates/` and are loaded via `//go:embed`. The old root `templates/` directory was deleted.

Same for frontend assets — `web/` is embedded via `web/embed.go`.

There is no runtime template hot-reload.

## Key model quirks

- `Port` has `validate:"min=0,max=65535"` (not `required`). Default is applied by `usecase.applyDefaults()`. Validation runs AFTER `Prepare()`.
- `ConfigRequest` fields are grouped by type (common / k8s / dockerfile) in struct comments at `internal/domain/types.go`.
- `K8sEnable` was removed — do not add it back.

## Generators

Three supported types: `compose`, `k8s`, `dockerfile`.

Adding a new generator type:
1. Implement `Generator` interface from `internal/infrastructure/generators/interface.go`
2. Register in `registry.go` map
3. Add templates under `internal/infrastructure/generators/templates/<type>/`

The K8s generator supports **20 resource kinds** (Deployment, Service, ConfigMap, Secret, …). Each kind has its own `.yaml.tmpl` file. The template cache uses sync.RWMutex with double-check locking — two goroutines do not parse the same template twice.

## Env var ordering

Compose and Dockerfile generators sort env keys alphabetically before rendering. Output is deterministic, not map-iteration-order dependent.

## DESIGN.md caveat

`DESIGN.md` lists Nginx, Docker Daemon, Containerd, ZIP packaging, template hot-reload, plugin system, and i18n as planned features. **None are implemented.** Only compose / k8s / dockerfile work.

## Test coverage

Only two packages have tests:
- `internal/infrastructure/generators/` — 89.5%
- `internal/infrastructure/storage/` — 74.5%

`internal/usecase/` and `internal/presentation/` have zero tests.

## Previous bugs fixed (do not reintroduce)

1. `applyDefaults` ran AFTER validation → moved `Prepare()` before `validate.Struct()`
2. `getConfig` returned 404 for all errors → now 404 only for `sql.ErrNoRows`, 500 otherwise
3. DSN assumed no `?` → now checks `strings.Contains(dsn, "?")`
4. Template cache TOCTOU → double-check inside write lock
5. Env map iteration non-deterministic → sorted keys
6. No graceful shutdown → signal handler with 5s timeout

## SQLite

- WAL journal mode enabled
- DSN separator handled: if DSN already has `?`, next param uses `&`
- Timestamps stored as RFC3339Nano text; `scanRecord` returns error on parse failure (no longer swallowed)
