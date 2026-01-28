# attgo_raw_string

**Priority:** MEDIUM (disabled by default)

## Description

Suggests using raw strings (backticks) over double-quoted strings with multiple escape sequences.

## Rationale

Raw strings improve readability when:
- String contains multiple quotes or backslashes
- String represents patterns, queries, or paths

## Examples

### Bad

```go
// Prometheus query with escaped quotes
query := "vouch_relay_execution_config_total{result=\"succeeded\"}"

// Windows path with many backslashes
path := "C:\\Users\\name\\Documents\\file.txt"

// JSON with escaped quotes
json := "{\"key\": \"value\", \"num\": 123}"
```

### Good

```go
// Prometheus query - much clearer
query := `vouch_relay_execution_config_total{result="succeeded"}`

// Windows path
path := `C:\Users\name\Documents\file.txt`

// JSON
json := `{"key": "value", "num": 123}`
```

### Exceptions (Not Flagged)

```go
// Only 1-2 escapes - below threshold
simple := "hello \"world\""

// Contains backtick - can't use raw string
code := "use `backticks` here"

// Newline/tab escapes serve a purpose
multiline := "line1\nline2"
tabbed := "col1\tcol2"
```

## Configuration

```yaml
settings:
  enable_raw_string: true  # Opt-in (disabled by default)
```

## Behavior

The rule triggers when:
- String has 3 or more escape sequences
- The escape sequences are `\"` or `\\` (not `\n`, `\t`, `\r`)
- The string doesn't contain backticks (which would make raw strings impossible)

## Suppression

```go
query := "intentionally \"escaped\"" //nolint:attgo_raw_string
```

## Source

- [attestant PR #722](https://github.com/attestantio/attestant/pull/722)
