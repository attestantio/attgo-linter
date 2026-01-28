# attgo-linter

Attestant Go Style Linter - enforces [attestantio](https://github.com/attestantio) organization coding standards as a [golangci-lint](https://golangci-lint.run/) module plugin.

## Important: How Module Plugins Work

**attgo-linter is a golangci-lint module plugin.** This means:

- You **cannot** use `golangci-lint run` directly - it will fail with "plugin not found"
- You **must** build a custom golangci-lint binary that includes this plugin
- The custom binary replaces `golangci-lint` for projects using attgo-linter

## Quick Start

### Step 1: Add Configuration Files

Add these two files to your project root:

**`.custom-gcl.yml`** (build configuration):
```yaml
version: v2.0.0
name: custom-gcl

plugins:
  - module: "github.com/attestantio/attgo-linter"
    version: v0.1.0
```

**`.golangci.yml`** (linter configuration):
```yaml
version: "2"

linters:
  enable:
    - attgo

  settings:
    custom:
      attgo:
        type: "module"
        description: "Attestant organization style linter"
        settings:
          enable_no_pkg_logger: true
          enable_enum_iota: true
          enable_current_year: true
```

### Step 2: Build the Custom Binary

```bash
# Requires golangci-lint v2.0+ installed
golangci-lint custom
```

This creates `./custom-gcl` (name from `.custom-gcl.yml`).

### Step 3: Run the Linter

```bash
# Use the custom binary instead of golangci-lint
./custom-gcl run
```

**Important:** Always use `./custom-gcl run`, not `golangci-lint run`.

## GitHub Actions Integration

Since module plugins require a custom binary, your CI workflow must build it first.

Add this workflow to `.github/workflows/golangci-lint.yml`:

```yaml
name: golangci-lint

on:
  push:
    branches: [master]
  pull_request:

permissions:
  contents: read

jobs:
  golangci:
    name: lint
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "^1.25"
          cache: false

      # Cache the custom binary to speed up subsequent runs
      - name: Cache custom golangci-lint
        id: cache-custom-gcl
        uses: actions/cache@v4
        with:
          path: ./custom-gcl
          key: custom-gcl-${{ hashFiles('.custom-gcl.yml') }}

      # Build custom binary only if not cached
      - name: Build custom golangci-lint
        if: steps.cache-custom-gcl.outputs.cache-hit != 'true'
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.0.0
          golangci-lint custom

      # Run using the custom binary
      - name: Run golangci-lint
        run: ./custom-gcl run --timeout=5m
```

## Local Development with Makefile

Add to your `Makefile`:

```makefile
CUSTOM_GCL := ./custom-gcl

# Build custom golangci-lint if needed
$(CUSTOM_GCL): .custom-gcl.yml
	golangci-lint custom

.PHONY: lint
lint: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run

.PHONY: lint-fix
lint-fix: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run --fix
```

Then run:
```bash
make lint
```

## Configuration Reference

### Full `.golangci.yml` Example

```yaml
version: "2"

linters:
  enable:
    - attgo

  settings:
    custom:
      attgo:
        type: "module"
        description: "Attestant organization style linter"
        settings:
          # HIGH PRIORITY - enabled by default
          enable_no_pkg_logger: true
          enable_enum_iota: true
          enable_current_year: true

          # MEDIUM PRIORITY - disabled by default
          enable_capital_comment: false
          enable_func_opts: false
          enable_raw_string: false

          # LOW PRIORITY - disabled by default
          enable_struct_field_order: false
          enable_interface_check: false

          # Custom logger patterns (optional)
          logger_type_patterns:
            - "zerolog.Logger"
            - "*zerolog.Logger"
            - "zap.Logger"
            - "*zap.Logger"

          # Custom enum suffixes (optional)
          enum_type_suffixes:
            - "Type"
            - "Status"
            - "State"
            - "Kind"
            - "Mode"
```

## Rules

### HIGH PRIORITY (Enabled by Default)

#### attgo_no_pkg_logger

Loggers must be struct fields, not package-level variables.

**Rationale:** Package-level loggers make testing difficult and prevent proper dependency injection. Struct field loggers enable:
- Injecting mock loggers for testing
- Clear ownership of logging configuration
- Better traceability in concurrent code

**Bad:**
```go
var log zerolog.Logger

func main() {
    log.Info().Msg("hello")
}
```

**Good:**
```go
type Service struct {
    log zerolog.Logger
}

func (s *Service) Run() {
    s.log.Info().Msg("hello")
}
```

**Configuration:**
```yaml
settings:
  logger_type_patterns:
    - "zerolog.Logger"
    - "*zerolog.Logger"
    - "zap.Logger"
    - "*zap.Logger"
```

---

#### attgo_enum_iota

Enum types should use `uint64` + iota pattern, not string constants.

**Rationale:** Integer-based enums provide several advantages:
- Memory efficiency (integers vs strings)
- Faster comparison operations
- Type safety preventing mixing of different enum types
- Exhaustive switch checking by static analyzers
- Explicit string representation via `String()` method

**Bad:**
```go
type SANType string

const (
    SANTypeDNS   SANType = "dns"
    SANTypeEmail SANType = "email"
)
```

**Good:**
```go
type SANType uint64

const (
    SANTypeUnknown SANType = iota
    SANTypeDNS
    SANTypeEmail
)

var sanTypeStrings = [...]string{"unknown", "dns", "email"}

func (s SANType) String() string {
    return sanTypeStrings[s]
}
```

**Configuration:**
```yaml
settings:
  enum_type_suffixes:
    - "Type"
    - "Status"
    - "State"
    - "Kind"
    - "Mode"
```

---

#### attgo_current_year

New files must have the current year in their copyright header.

**Rationale:** Accurate copyright years are important for:
- Legal compliance
- Indicating when code was created/modified
- Consistency across the codebase

**Bad (in 2026):**
```go
// Copyright © 2025 Attestant Limited.
```

**Good:**
```go
// Copyright © 2026 Attestant Limited.
```

**Also acceptable (year ranges):**
```go
// Copyright © 2023-2026 Attestant Limited.
```

---

### MEDIUM PRIORITY (Disabled by Default)

#### attgo_capital_comment

Comments should start with a capital letter.

**Rationale:** Consistency and readability. Well-formatted comments indicate attention to code quality.

**Exceptions:**
- Comments starting with identifiers (`someFunc is...`)
- nolint directives
- URLs, TODOs, build tags
- License boilerplate text

**Bad:**
```go
// this is a comment
```

**Good:**
```go
// This is a comment
// someVariable contains the value
```

---

#### attgo_func_opts

Service constructors with many parameters should use the functional options pattern.

**Rationale:** The functional options pattern provides:
- Extensibility without breaking existing callers
- Self-documenting named options
- Natural handling of optional parameters
- Easy default values
- Validation within option functions

**Bad:**
```go
func NewService(log Logger, db DB, cache Cache, timeout time.Duration) *Service
```

**Good:**
```go
type Option func(*Service)

func WithLogger(log Logger) Option {
    return func(s *Service) { s.log = log }
}

func New(opts ...Option) *Service
```

---

#### attgo_raw_string

Prefer raw strings (backticks) over heavily escaped double-quoted strings.

**Rationale:** Raw strings improve readability when the string contains multiple quotes or backslashes, such as in queries, paths, or JSON.

**Bad:**
```go
query := "vouch_relay_execution_config_total{result=\"succeeded\"}"
path := "C:\\Users\\name\\Documents\\file.txt"
```

**Good:**
```go
query := `vouch_relay_execution_config_total{result="succeeded"}`
path := `C:\Users\name\Documents\file.txt`
```

---

### LOW PRIORITY (Disabled by Default)

#### attgo_struct_field_order

Struct fields should be ordered: logger → metrics → dependencies → data → sync.

**Rationale:** Consistent field ordering creates predictable structure:
- Know where to find fields without searching
- Easier code review
- Faster onboarding for new developers
- Logical grouping of related fields

**Example:**
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
    cache  map[string]Value

    // Synchronization
    mu   sync.Mutex
    done chan struct{}
}
```

---

#### attgo_interface_check

Suggests adding `var _ Interface = (*Struct)(nil)` compile-time checks.

**Rationale:** This pattern provides:
- Compile-time verification that a struct implements an interface
- Clear documentation of which interfaces a type implements
- Immediate build failure if methods are removed or signatures change
- Better IDE support for code navigation

**Pattern:**
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type MyReader struct{}

// Compile-time check
var _ Reader = (*MyReader)(nil)

func (r *MyReader) Read(p []byte) (int, error) {
    return 0, nil
}
```

---

## Disabling Rules

Use standard golangci-lint nolint directives:

```go
var log zerolog.Logger //nolint:attgo_no_pkg_logger // legacy code

query := "escaped\"string" //nolint:attgo_raw_string // intentional
```

## Troubleshooting

### "plugin 'attgo' not found"

You're using `golangci-lint run` instead of the custom binary. Use:
```bash
./custom-gcl run  # NOT golangci-lint run
```

### Build fails with module errors

Clear the cache and rebuild:
```bash
rm -rf ~/.cache/golangci-lint
rm ./custom-gcl
golangci-lint custom
```

### Plugin version not found

Ensure the version in `.custom-gcl.yml` matches a published release:
```yaml
plugins:
  - module: "github.com/attestantio/attgo-linter"
    version: v0.1.0  # Must be a valid git tag
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
