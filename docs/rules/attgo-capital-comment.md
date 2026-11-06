# attgo_capital_comment

**Priority:** MEDIUM (disabled by default)

## Description

Checks that comments start with a capital letter.

## Rationale

- **Consistency**: Uniform style across the codebase
- **Readability**: Properly formatted comments are easier to read
- **Professionalism**: Well-formatted comments indicate care for code quality

## Examples

### Bad

```go
// this function does something
func DoSomething() {}

// returns the value
func GetValue() int { return 0 }
```

### Good

```go
// This function does something
func DoSomething() {}

// GetValue returns the value
func GetValue() int { return 0 }
```

### Allowed Exceptions

```go
// someFunc is used to process data (identifier reference)
// myVariable contains the configuration (camelCase identifier)
// my_var stores the value (snake_case identifier)
// nolint:errcheck (directive)
// TODO: fix this later (TODO marker)
// See https://example.com (URL)
// +build linux (build tag)
// go:generate stringer -type=Foo (go directive)
// ... continued from above (punctuation)
// 123 is the magic number (number)
```

## Configuration

```yaml
settings:
  enable_capital_comment: true  # Opt-in (disabled by default)
```

## Suppression

The directive names the linter, `attgo`, not the rule, so it silences every `attgo` rule on the comment. Put it on its own line below the comment text, in the same comment block. On a doc comment it also covers the declaration. See [Linter Names](../../README.md#linter-names).

```go
// this is intentionally lowercase
//nolint:attgo
```

## Notes

- Identifier references are detected by looking for patterns like `someFunc is...`, `myVar contains...`
- Common English words like "this", "see", "use" are not treated as identifiers
- The rule aims to catch genuine style violations while avoiding false positives on technical comments

## Source

- [vouch PR #334](https://github.com/attestantio/vouch/pull/334)
