# Maintenance Guide

This document describes how to add, modify, and maintain rules in attgo-linter.

## Adding a New Rule

### 1. Create the Analyzer Package

Create a new directory under `analyzers/`:

```
analyzers/
└── newrule/
    ├── analyzer.go
    └── analyzer_test.go
```

### 2. Implement the Analyzer

```go
// analyzers/newrule/analyzer.go
package newrule

import (
    "golang.org/x/tools/go/analysis"
)

const (
    analyzerName = "attgo_new_rule"  // Use underscores, not hyphens
    doc          = `description of what the rule checks

Detailed explanation...

Bad:
    example of bad code

Good:
    example of good code`
)

// Analyzer is the new rule analyzer.
var Analyzer = &analysis.Analyzer{
    Name: analyzerName,
    Doc:  doc,
    Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
    // Implementation
    return nil, nil
}
```

### 3. Add Test Fixtures

Create test data:

```
analyzers/newrule/testdata/src/newrule/
└── newrule.go
```

Use `// want` comments for expected diagnostics:

```go
// testdata/src/newrule/newrule.go
package newrule

var badCode = something // want `expected error message`

var goodCode = something // No warning
```

### 4. Add the Test

```go
// analyzer_test.go
package newrule_test

import (
    "testing"

    "github.com/attestantio/attgo-linter/analyzers/newrule"
    "golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, newrule.Analyzer, "newrule")
}
```

### 5. Register in Plugin

Update `plugin.go`:

```go
import (
    // ...
    "github.com/attestantio/attgo-linter/analyzers/newrule"
)

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
    // ...
    if p.cfg.EnableNewRule {
        analyzers = append(analyzers, newrule.Analyzer)
    }
    // ...
}
```

### 6. Add Configuration

Update `config.go`:

```go
type Config struct {
    // ...
    EnableNewRule bool `json:"enable_new_rule"`
}

func DefaultConfig() *Config {
    return &Config{
        // Set default (true for HIGH priority, false for MEDIUM/LOW)
        EnableNewRule: false,
    }
}
```

Update `plugin.go` to handle the config:

```go
if _, ok := rawSettings["enable_new_rule"]; ok {
    cfg.EnableNewRule = userCfg.EnableNewRule
}
```

### 7. Document the Rule

1. Update `README.md` with rule documentation
2. Create `docs/rules/attgo-new-rule.md` for detailed documentation
3. Update `CHANGELOG.md`

## Rule Naming Conventions

- Analyzer names use underscores: `attgo_new_rule`
- Package names use lowercase without separators: `newrule`
- Config fields use snake_case: `enable_new_rule`

## Testing

Run all tests:
```bash
go test ./...
```

Run specific analyzer tests:
```bash
go test ./analyzers/newrule/...
```

Run with verbose output:
```bash
go test -v ./analyzers/newrule/...
```

## Common Patterns

### Type-Aware Analysis

For analyzers that need type information, the plugin uses `LoadModeTypesInfo`:

```go
func (p *Plugin) GetLoadMode() string {
    return register.LoadModeTypesInfo
}
```

This gives access to `pass.TypesInfo` in your analyzer.

### Pattern Matching Types

```go
func isTargetType(t types.Type, pattern string) bool {
    typeName := types.TypeString(t, nil)
    return strings.HasSuffix(typeName, pattern)
}
```

### Reporting Diagnostics

```go
pass.Reportf(node.Pos(), "message with %s formatting", arg)
```

### Suggested Fixes (Optional)

```go
pass.Report(analysis.Diagnostic{
    Pos:     node.Pos(),
    Message: "description",
    SuggestedFixes: []analysis.SuggestedFix{
        {
            Message: "fix description",
            TextEdits: []analysis.TextEdit{
                {
                    Pos:     start,
                    End:     end,
                    NewText: []byte("replacement"),
                },
            },
        },
    },
})
```

## Release Process

1. Update `CHANGELOG.md` with changes
2. Create a git tag: `git tag v0.x.0`
3. Push: `git push origin v0.x.0`
4. Users update their `.custom-gcl.yml` to use the new version
