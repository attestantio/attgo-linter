# attgo_interface_check

**Priority:** LOW (disabled by default)

## Description

Suggests adding compile-time interface compliance checks when a struct implements an interface.

## Rationale

The pattern `var _ Interface = (*Struct)(nil)` provides:

1. **Compile-time Safety**: Missing methods cause build failures, not runtime errors
2. **Documentation**: Explicitly shows which interfaces a type implements
3. **Refactoring Safety**: Adding/removing methods from interfaces immediately shows affected types
4. **IDE Support**: Better code navigation and autocomplete

## Examples

### Without Check (Risky)

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type MyReader struct{}

func (r *MyReader) Read(p []byte) (int, error) {
    return 0, nil
}

// If someone changes the interface or removes the method,
// the error only appears at runtime when Read is called.
```

### With Check (Safe)

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type MyReader struct{}

// Compile-time check - fails immediately if MyReader
// doesn't implement Reader
var _ Reader = (*MyReader)(nil)

func (r *MyReader) Read(p []byte) (int, error) {
    return 0, nil
}
```

## Configuration

```yaml
settings:
  enable_interface_check: true  # Opt-in (disabled by default)
```

## Behavior

The rule:
1. Finds all interfaces with methods defined in the package
2. Finds all struct types in the package
3. Checks if each struct implements any interface (via pointer or value receiver)
4. Reports if there's no `var _ Interface = (*Struct)(nil)` check

## Suppression

```go
type MyReader struct{} //nolint:attgo_interface_check
```

## Notes

- Only checks interfaces defined in the same package
- Empty interfaces (no methods) are ignored
- Both value and pointer receivers are considered
- Existing checks with the correct pattern are recognized and not flagged
