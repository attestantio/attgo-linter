# attgo_struct_field_order

**Priority:** LOW (disabled by default)

## Description

Enforces a consistent ordering of struct fields by category: logger, metrics, dependencies, data, synchronization.

## Rationale

Consistent field ordering:
1. **Predictability**: Know where to find fields without searching
2. **Code Review**: Easier to review when structure is consistent
3. **Onboarding**: New developers learn the pattern once
4. **Logical Grouping**: Related fields are together

## Field Categories

1. **Logger** - Logging fields (log, logger)
2. **Metrics** - Monitoring fields (metrics, monitor)
3. **Dependencies** - External services (client, db, cache, service)
4. **Data** - Configuration and state (config, name, value)
5. **Synchronization** - Concurrency primitives (mutex, wg, channels)

## Examples

### Bad

```go
type Service struct {
    mu     sync.Mutex      // Sync should be last
    log    zerolog.Logger  // Logger should be first
    config Config
    db     Database        // Dependency should come before data
}
```

### Good

```go
type Service struct {
    // Logger
    log zerolog.Logger

    // Metrics
    metrics *prometheus.Registry

    // Dependencies
    client *http.Client
    db     Database

    // Data
    config Config
    name   string

    // Synchronization
    mu   sync.Mutex
    done chan struct{}
}
```

## Configuration

```yaml
settings:
  enable_struct_field_order: true  # Opt-in (disabled by default)
```

## Detection Rules

Fields are categorized by name and type:

| Category | Detected By |
|----------|-------------|
| Logger | Names: `log`, `logger`, `*log`, `*logger` |
| Metrics | Names: `metrics`, `monitor`, `*metrics` |
| Dependency | Names ending in: `client`, `service`, `provider`, `handler`, `store`, `repo` |
| Sync | Types: `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, channels |
| Data | Everything else |

## Suppression

```go
type Service struct { //nolint:attgo_struct_field_order
    // Custom ordering for specific reason
    mu  sync.Mutex
    log zerolog.Logger
}
```

## Notes

- Embedded fields are ignored
- Only named fields are checked
- The rule reports when a field from an earlier category appears after a field from a later category
