# Integration Guide

This guide covers integrating attgo-linter into your project and CI/CD pipeline.

## Understanding Module Plugins

**Important:** attgo-linter is a golangci-lint **module plugin**. This means:

1. You **cannot** use `golangci-lint run` directly
2. You **must** build a custom binary using `golangci-lint custom`
3. You then use that custom binary instead of `golangci-lint`

This is how all golangci-lint module plugins work - they must be compiled into the binary at build time.

## Prerequisites

- Go 1.23 or later
- golangci-lint v2.0.0 or later

Install golangci-lint:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.0.0
```

## Step-by-Step Setup

### 1. Create `.custom-gcl.yml`

This file tells `golangci-lint custom` which plugins to compile:

```yaml
version: v2.0.0
name: custom-gcl

plugins:
  - module: "github.com/attestantio/attgo-linter"
    version: v0.1.0
```

**Note:** The `name` field determines the output binary name (`custom-gcl`).

### 2. Create `.golangci.yml`

Configure which rules to enable:

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
          # Enable the rules you want
          enable_no_pkg_logger: true
          enable_enum_iota: true
          enable_current_year: true
```

### 3. Build the Custom Binary

```bash
golangci-lint custom
```

This creates `./custom-gcl` in your current directory.

### 4. Run the Linter

```bash
./custom-gcl run
```

**Always use `./custom-gcl run`, not `golangci-lint run`.**

## GitHub Actions

### Basic Workflow

```yaml
name: golangci-lint

on:
  push:
    branches: [master, main]
  pull_request:

permissions:
  contents: read

jobs:
  lint:
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "^1.23"
          cache: false

      # Cache the custom binary
      - name: Cache custom golangci-lint
        id: cache
        uses: actions/cache@v4
        with:
          path: ./custom-gcl
          key: custom-gcl-${{ hashFiles('.custom-gcl.yml') }}

      # Build only if not cached
      - name: Build custom golangci-lint
        if: steps.cache.outputs.cache-hit != 'true'
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.0.0
          golangci-lint custom

      # Run the linter
      - name: Lint
        run: ./custom-gcl run --timeout=5m
```

### With New Issues Only

To only report issues on changed lines (useful for existing codebases):

```yaml
      - name: Lint
        run: ./custom-gcl run --timeout=5m --new-from-rev=origin/${{ github.base_ref }}
        if: github.event_name == 'pull_request'

      - name: Lint (full)
        run: ./custom-gcl run --timeout=5m
        if: github.event_name != 'pull_request'
```

## GitLab CI

```yaml
lint:
  stage: test
  image: golang:1.23
  cache:
    paths:
      - ./custom-gcl
    key: custom-gcl-${CI_COMMIT_REF_SLUG}
  script:
    - |
      if [ ! -f ./custom-gcl ]; then
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.0.0
        golangci-lint custom
      fi
    - ./custom-gcl run --timeout=5m
```

## Makefile Integration

```makefile
CUSTOM_GCL := ./custom-gcl

# Build custom binary if .custom-gcl.yml is newer or binary doesn't exist
$(CUSTOM_GCL): .custom-gcl.yml
	golangci-lint custom

.PHONY: lint
lint: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run

.PHONY: lint-fix
lint-fix: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run --fix

.PHONY: lint-new
lint-new: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run --new-from-rev=HEAD~1

.PHONY: clean-lint
clean-lint:
	rm -f $(CUSTOM_GCL)
```

## Pre-commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
set -e

# Build custom binary if needed
if [ ! -f ./custom-gcl ] || [ .custom-gcl.yml -nt ./custom-gcl ]; then
    echo "Building custom golangci-lint..."
    golangci-lint custom
fi

# Run on staged files only
./custom-gcl run --new-from-rev=HEAD
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Local Development for attgo-linter

When developing attgo-linter itself, use `path` instead of `version`:

```yaml
# .custom-gcl.yml for local development
version: v2.0.0
name: custom-gcl-attgo

plugins:
  - module: "github.com/attestantio/attgo-linter"
    path: "."  # Points to local directory
```

## Troubleshooting

### "plugin 'attgo' not found"

**Cause:** You're using `golangci-lint run` instead of the custom binary.

**Fix:** Use the custom binary:
```bash
./custom-gcl run  # Correct
# NOT: golangci-lint run
```

### "module not found" during build

**Cause:** The version in `.custom-gcl.yml` doesn't exist.

**Fix:** Use a valid release tag:
```yaml
plugins:
  - module: "github.com/attestantio/attgo-linter"
    version: v0.1.0  # Must be a valid git tag
```

### Build fails with cache issues

**Fix:** Clear caches and rebuild:
```bash
rm -rf ~/.cache/golangci-lint
rm ./custom-gcl
golangci-lint custom
```

### CI caching not working

Ensure your cache key includes the config file hash:
```yaml
key: custom-gcl-${{ hashFiles('.custom-gcl.yml') }}
```

## Version Compatibility

| attgo-linter | golangci-lint | Go |
|--------------|---------------|-----|
| v0.1.x | v2.0.0+ | 1.23+ |
