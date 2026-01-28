# attgo_func_opts

**Priority:** MEDIUM (disabled by default)

## Description

Suggests using the functional options pattern for service type constructors that have many parameters.

## Rationale

The functional options pattern provides:

1. **Extensibility**: Add new options without breaking existing callers
2. **Readability**: Named options are self-documenting
3. **Optional Parameters**: Natural way to handle optional configuration
4. **Default Values**: Easy to provide sensible defaults
5. **Validation**: Options can validate their values

## Examples

### Bad

```go
type UserService struct {
    db      Database
    cache   Cache
    logger  Logger
    metrics Metrics
}

func NewUserService(db Database, cache Cache, logger Logger, metrics Metrics) *UserService {
    return &UserService{
        db:      db,
        cache:   cache,
        logger:  logger,
        metrics: metrics,
    }
}
```

### Good

```go
type UserService struct {
    db      Database
    cache   Cache
    logger  Logger
    metrics Metrics
}

type Option func(*UserService)

func WithDatabase(db Database) Option {
    return func(s *UserService) {
        s.db = db
    }
}

func WithCache(cache Cache) Option {
    return func(s *UserService) {
        s.cache = cache
    }
}

func WithLogger(logger Logger) Option {
    return func(s *UserService) {
        s.logger = logger
    }
}

func WithMetrics(metrics Metrics) Option {
    return func(s *UserService) {
        s.metrics = metrics
    }
}

func New(opts ...Option) *UserService {
    s := &UserService{
        // Default values
        logger: zerolog.Nop(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## Configuration

```yaml
settings:
  enable_func_opts: true  # Opt-in (disabled by default)
```

## Behavior

The rule triggers when:
- A function is named `New...` or `Create...`
- It returns a pointer to a service-like type (suffix: Service, Manager, Handler, Controller, Provider, Client, Server)
- It has more than 3 non-context parameters
- It doesn't already use variadic options (e.g., `...Option`)

## Suppression

```go
func NewService(a, b, c, d Interface) *Service { //nolint:attgo_func_opts
    // ...
}
```

## Reference

See [vouch/services/attester/standard/parameters.go](https://github.com/attestantio/vouch) for the canonical implementation pattern.
