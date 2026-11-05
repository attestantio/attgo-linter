# Changelog

## dev

### Added

- `RegisterAnalyzer` registers an analyzer as a golangci-lint plugin of its own, named after the analyzer. A repository enables it with a `linters.settings.custom` entry of that name, with no `settings` (plus a `linters.enable` entry when it sets `linters.default` to `none` or `fast`), and suppresses it with `//nolint:<name>`, which leaves every other linter reporting. The existing rules stay in the `attgo` plugin, and existing `attgo` configurations need no change.

### Fixed

- README: `//nolint:attgo_no_pkg_logger` and `//nolint:attgo_raw_string` suppressed nothing, because golangci-lint reports every rule of the plugin under the linter name `attgo`. The examples now use `//nolint:attgo`.
- `docs/rules`: the suppression examples use `//nolint:attgo`, and each now holds a violation that its rule reports and the directive silences. The source links of `attgo_current_year` and `attgo_raw_string` point to public repositories.

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
