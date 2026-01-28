# Changelog

## dev

- (work in progress)

## v0.1.0

Initial release.

### Added

- `attgo-no-pkg-logger` rule: loggers must be struct fields, not package-level variables
- `attgo-enum-iota` rule: enum types should use `uint64` + iota pattern
- `attgo-current-year` rule: new files must have current year in copyright header
- `attgo-capital-comment` rule: comments should start with capital letter
- `attgo-func-opts` rule: services should use functional options pattern
- `attgo-raw-string` rule: prefer raw strings over escaped strings
- `attgo-struct-field-order` rule: struct field ordering convention
- `attgo-interface-check` rule: suggest interface compliance checks
- golangci-lint module plugin integration
- Agent (e.g. Claude Code) skill for rule management
