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

### 5. Register Under Its Own Name

golangci-lint reports each issue under the name of the plugin that produced it, and `linters.enable` and `//nolint` match that name. Every rule inside the `attgo` plugin shares the name `attgo`, so `//nolint:attgo` silences them all. A new rule therefore registers as a plugin of its own, named after the analyzer.

Add one call to the `init` function in `plugin.go`:

```go
import (
    // ...
    "github.com/attestantio/attgo-linter/analyzers/newrule"
)

func init() {
    register.Plugin("attgo", New)

    RegisterAnalyzer(newrule.Analyzer)
}
```

That is all the code. There is no `Config` field and no `enable_*` setting. A repository turns the rule on by adding a `linters.settings.custom` entry named after the analyzer. A config that sets `linters.default` to `none` or `fast` must also list the name under `linters.enable`, as here:

```yaml
linters:
  default: none
  enable:
    - attgo_new_rule

  settings:
    custom:
      attgo_new_rule:
        type: "module"
        description: "What the rule checks"
```

The entry takes no `settings`; `RegisterAnalyzer` rejects any. The rule registers from `init` in the root package, which a consumer's attgo-linter entry in `.custom-gcl.yml` already imports, so that file needs no extra entry, only a `version` that includes the rule.

The rule is then suppressed with `//nolint:attgo_new_rule`, which leaves every other linter reporting on that line.

The `enable_*` rules predate this. The `attgo` plugin builds them (`Config`, `New` and `BuildAnalyzers`), and they stay there until they move to their own names. Do not add new rules there.

An analyzer value belongs to one plugin. golangci-lint maps each analyzer value to one linter name, so when the `attgo` plugin and a per-name plugin both return the same value, its issues report under only one of the two names, and `//nolint:attgo` stops matching them. Moving an existing rule to its own name means removing it from the `attgo` plugin in the same change.

### 6. Document the Rule

1. Update `README.md` with rule documentation, including the `linters.settings.custom` entry
2. Create `docs/rules/attgo-new-rule.md` for detailed documentation
3. Update `CHANGELOG.md`

## Rule Naming Conventions

- Analyzer names use underscores and are lowercase: `attgo_new_rule`. The analyzer name is also the linter name, and golangci-lint lowercases the names in `//nolint` directives, so a name with capitals can never be suppressed.
- Package names use lowercase without separators: `newrule`
- Settings of the `attgo` plugin use snake_case keys, such as `enable_no_pkg_logger`

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

### Per-Name Registration

`TestIntegration_PerNameNolint` builds a custom golangci-lint binary from the checkout, the way a consumer does, and runs it over `testdata/project`. Two sample analyzers from `testdata/sampleplugin` register under their own names next to the `attgo` plugin. The test checks that a `//nolint` directive naming one of them silences only that one. It runs once with the `attgo` settings of vouch (`vouch.yml`, the same as dirk's) and once with those of go-eth2-client (`go-eth2-client.yml`). A lowercase comment in `testdata/project` is reported under the first and not under the second, because go-eth2-client turns `enable_capital_comment` off, so the test also shows that those settings still apply next to per-name plugins.

The test builds the golangci-lint version of `.custom-gcl.yml`. It needs `golangci-lint`, `git` and network access on the machine, takes about a minute, and `-short` skips it. When `CI` is set and `golangci-lint` is missing, it fails instead of skipping. Run it with `-count=1`, because `go test` does not see an edit to `testdata/sampleplugin` and would report a cached pass:

```bash
go test -count=1 -run TestIntegration_PerNameNolint -v .
```

Keep `testdata/project` free of imports: an older golangci-lint cannot always read the export data of a newer Go standard library.

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
