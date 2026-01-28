# Add Rule Skill

This skill helps you add new rules to the attgo-linter project.

## Usage

When asked to add a new linter rule, follow these steps:

### 1. Gather Information

Ask the user for:
- **Rule name**: e.g., "no-global-vars"
- **Description**: What the rule checks
- **Priority**: HIGH (enabled by default), MEDIUM, or LOW (disabled by default)
- **Bad example**: Code that should trigger the warning
- **Good example**: Code that passes the check
- **Configuration options**: Any customizable settings

### 2. Create Analyzer Package

Create the following files:

```
analyzers/{rulename}/
├── analyzer.go
└── analyzer_test.go
```

The analyzer name must use underscores (not hyphens): `attgo_{rule_name}`

### 3. Create Test Fixtures

```
analyzers/{rulename}/testdata/src/{rulename}/
└── {rulename}.go
```

Use `// want \`regex pattern\`` comments for expected diagnostics.

### 4. Register in Plugin

Update `plugin.go`:
- Add import
- Add to `BuildAnalyzers()` with config check

Update `config.go`:
- Add `Enable{RuleName}` field to `Config` struct
- Set default in `DefaultConfig()`

Update `plugin.go` config handling:
- Add explicit check for the setting in `New()`

### 5. Update Documentation

- Add rule section to `README.md`
- Create `docs/rules/attgo-{rule-name}.md`
- Update `CHANGELOG.md`

### 6. Verify

Run tests:
```bash
go test ./...
```

## Template: analyzer.go

```go
// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// ...

// Package {rulename} provides an analyzer that {description}.
package {rulename}

import (
    "golang.org/x/tools/go/analysis"
)

const (
    analyzerName = "attgo_{rule_name}"
    doc          = `{short description}

{detailed explanation}

Bad:
    {bad example}

Good:
    {good example}`
)

// Analyzer is the {description} analyzer.
var Analyzer = &analysis.Analyzer{
    Name: analyzerName,
    Doc:  doc,
    Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
    for _, file := range pass.Files {
        // Implementation
    }
    return nil, nil
}
```

## Template: analyzer_test.go

```go
// Copyright © 2025 Attestant Limited.
// ...

package {rulename}_test

import (
    "testing"

    "github.com/attestantio/attgo-linter/analyzers/{rulename}"
    "golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, {rulename}.Analyzer, "{rulename}")
}
```

## Template: testdata

```go
// Copyright © 2025 Attestant Limited.
// ...

package {rulename}

// Bad case.
var bad = something // want `expected error`

// Good case.
var good = something // No warning
```

## Checklist

- [ ] Created `analyzers/{rulename}/analyzer.go`
- [ ] Created `analyzers/{rulename}/analyzer_test.go`
- [ ] Created `analyzers/{rulename}/testdata/src/{rulename}/{rulename}.go`
- [ ] Added import to `plugin.go`
- [ ] Added to `BuildAnalyzers()` in `plugin.go`
- [ ] Added config field to `config.go`
- [ ] Added default value in `DefaultConfig()`
- [ ] Added config handling in `New()`
- [ ] Tests pass: `go test ./...`
- [ ] Added documentation to `README.md`
- [ ] Updated `CHANGELOG.md`
