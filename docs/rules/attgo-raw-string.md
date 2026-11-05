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

The directive names the linter, `attgo`, not the rule, so it silences every `attgo` rule on its line. See [Linter Names](../../README.md#linter-names).

```go
query := "{\"key\": \"value\"}" //nolint:attgo
```

## Source

- [go-eth2-client `api/v1/validatorstate.go`](https://github.com/attestantio/go-eth2-client/blob/6d6c8310d52ef56f6382065a338a2e0bc34e5d39/api/v1/validatorstate.go#L78-L96) and [vouch `services/blockrelay/configversion.go`](https://github.com/attestantio/vouch/blob/83758284525e3abb5bf46a1680a944020dbf831f/services/blockrelay/configversion.go#L45-L47) use raw strings for quoted text
- Enabled in [vouch](https://github.com/attestantio/vouch/blob/83758284525e3abb5bf46a1680a944020dbf831f/.golangci.yml#L207), [go-eth2-client](https://github.com/attestantio/go-eth2-client/blob/6d6c8310d52ef56f6382065a338a2e0bc34e5d39/.golangci.yml#L163) and [dirk](https://github.com/attestantio/dirk/blob/0e84fe34850ce9f9acd37069dcf8fc5522e309d3/.golangci.yml#L208)
