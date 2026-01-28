# attgo_current_year

**Priority:** HIGH (enabled by default)

## Description

Checks that copyright headers in Go files contain the current year.

## Rationale

Accurate copyright years are important for:

1. **Legal Compliance**: Copyright notices should reflect when the work was created/modified
2. **Freshness Indication**: Helps identify recently maintained code
3. **Consistency**: All files in a project should follow the same convention

## Examples

### Bad (in 2025)

```go
// Copyright © 2023 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

### Good

```go
// Copyright © 2025 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

### Also Acceptable (year ranges)

```go
// Copyright © 2023-2025 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

## Configuration

```yaml
settings:
  enable_current_year: true
```

## Suppression

```go
// Copyright © 2020 Attestant Limited. //nolint:attgo_current_year
```

## Notes

- This rule only checks the year in the copyright header, not the full format (use `goheader` linter for format validation)
- Year ranges like "2023-2025" are valid if the end year is current
- Files without copyright headers are not flagged (that's a separate concern)

## Source

- [attestant PR #719](https://github.com/attestantio/attestant/pull/719)
