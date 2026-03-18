# attgo_current_year

**Priority:** HIGH (enabled by default)

## Description

Checks that copyright headers in new or modified Go files contain the current year. Uses git to determine which files have changed relative to a configurable base ref (default: `auto`, which detects the remote's default branch), so untouched files with older copyright years are not flagged. Uses three-dot diff (`base...HEAD`) to compare against the merge base, so only changes on the current branch are detected.

## Behavior

- **New files** (added in git): must have the current year. Flags with a message suggesting the current year.
- **Modified files** (changed in git): must have the current year as the end of a range. Flags with a message suggesting `<original year>-<current year>`.
- **Unchanged files**: not checked — old copyright years are fine if the file was not touched.
- **Git unavailable**: falls back to checking all files (backward compatible).

## Rationale

Accurate copyright years are important for:

1. **Legal Compliance**: Copyright notices should reflect when the work was created/modified
2. **Freshness Indication**: Helps identify recently maintained code
3. **Consistency**: All files in a project should follow the same convention

## Examples

### New file — Bad

```go
// Copyright © 2023 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

### New file — Good

```go
// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

### Modified file — Bad

```go
// Copyright © 2023 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

### Modified file — Good

```go
// Copyright © 2023-2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0
```

### Untouched file — Not flagged

```go
// Copyright © 2020 Attestant Limited.
// (file has not been modified — this is fine)
```

## Configuration

```yaml
settings:
  enable_current_year: true
  # Default: "auto" (detects the remote's default branch).
  # Set to an explicit ref like "origin/main" or "origin/master" to override.
  # Set to "" to check all files (disables git-aware filtering).
  current_year_base_ref: "auto"
```

## Suppression

```go
// Copyright © 2020 Attestant Limited. //nolint:attgo_current_year
```

## Notes

- This rule only checks the year in the copyright header, not the full format (use `goheader` linter for format validation)
- Year ranges like "2023-2026" are valid if the end year is current
- Files without copyright headers are not flagged (that's a separate concern)
- The `current_year_base_ref` setting controls which git ref to diff against; `"auto"` detects the remote default branch, or set an explicit ref
- Renamed files are treated as modified (the copyright in the new path is checked)

## Source

- [attestant PR #719](https://github.com/attestantio/attestant/pull/719)
