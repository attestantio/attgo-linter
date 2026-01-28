# attgo_enum_iota

**Priority:** HIGH (enabled by default)

## Description

Enforces the iota pattern for enum types. Enum types should use `uint64` (or another integer type) with iota, not string constants.

## Rationale

Integer-based enums provide several advantages:

1. **Memory Efficiency**: Integers use less memory than strings
2. **Performance**: Integer comparison is faster than string comparison
3. **Type Safety**: Named types prevent mixing different enums
4. **Exhaustive Checking**: Static analyzers can verify switch statements cover all cases
5. **Serialization Control**: String representation is explicit via `String()` method

## Examples

### Bad

```go
type SANType string

const (
    SANTypeDNS   SANType = "dns"
    SANTypeEmail SANType = "email"
    SANTypeIP    SANType = "ip"
)
```

### Good

```go
type SANType uint64

const (
    SANTypeUnknown SANType = iota
    SANTypeDNS
    SANTypeEmail
    SANTypeIP
)

var sanTypeStrings = [...]string{
    "unknown",
    "dns",
    "email",
    "ip",
}

func (s SANType) String() string {
    if int(s) >= len(sanTypeStrings) {
        return "unknown"
    }
    return sanTypeStrings[s]
}

func SANTypeFromString(s string) SANType {
    for i, str := range sanTypeStrings {
        if str == s {
            return SANType(i)
        }
    }
    return SANTypeUnknown
}
```

## Configuration

```yaml
settings:
  enable_enum_iota: true
  enum_type_suffixes:
    - "Type"
    - "Status"
    - "State"
    - "Kind"
    - "Mode"
```

## Suppression

```go
type Color string

const (
    ColorRed Color = "red" //nolint:attgo_enum_iota // intentionally string-based
)
```

## Reference Implementation

See [go-eth2-client/spec/dataversion.go](https://github.com/attestantio/go-eth2-client/blob/master/spec/dataversion.go) for the canonical implementation pattern.

## Source

- [go-certmanager PR #1](https://github.com/attestantio/go-certmanager/pull/1)
