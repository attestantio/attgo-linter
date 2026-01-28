# attgo_no_pkg_logger

**Priority:** HIGH (enabled by default)

## Description

Detects package-level logger variables. Loggers should be struct fields, not package-level variables.

## Rationale

Package-level loggers create several problems:

1. **Difficult Testing**: Cannot inject mock loggers for testing
2. **Hidden Dependencies**: Logger dependency is implicit rather than explicit
3. **Initialization Order**: Package-level loggers may not be properly configured when used
4. **No Context**: Cannot associate logs with specific service instances

## Examples

### Bad

```go
package myservice

import "github.com/rs/zerolog"

var log zerolog.Logger // Package-level logger

func DoSomething() {
    log.Info().Msg("doing something")
}
```

### Good

```go
package myservice

import "github.com/rs/zerolog"

type Service struct {
    log zerolog.Logger // Struct field logger
}

func New(log zerolog.Logger) *Service {
    return &Service{log: log}
}

func (s *Service) DoSomething() {
    s.log.Info().Msg("doing something")
}
```

## Configuration

```yaml
settings:
  enable_no_pkg_logger: true
  logger_type_patterns:
    - "zerolog.Logger"
    - "*zerolog.Logger"
    - "zap.Logger"
    - "*zap.Logger"
    - "zap.SugaredLogger"
    - "*zap.SugaredLogger"
    - "logrus.Logger"
    - "*logrus.Logger"
    - "slog.Logger"
    - "*slog.Logger"
```

## Suppression

```go
var log zerolog.Logger //nolint:attgo_no_pkg_logger // legacy code, will refactor
```

## Source

- [go-certmanager PR #1](https://github.com/attestantio/go-certmanager/pull/1)
