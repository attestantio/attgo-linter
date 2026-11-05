package project

// Every literal below is reported by attgo_raw_string (six escape sequences)
// and by both sample analyzers. The directive on its line decides which of
// the three linters stay silent.
//
// This file has no copyright header, which attgo_current_year would flag once
// its year is stale. Keep it free of imports, see MAINTENANCE.md.
var (
	plain = "flag-me \"a\" \"b\" \"c\""

	byFirst = "flag-me \"a\" \"b\" \"c\"" //nolint:sample_first

	bySecond = "flag-me \"a\" \"b\" \"c\"" //nolint:sample_second

	byAttgo = "flag-me \"a\" \"b\" \"c\"" //nolint:attgo

	byBoth = "flag-me \"a\" \"b\" \"c\"" //nolint:sample_first,sample_second

	byRuleName = "flag-me \"a\" \"b\" \"c\"" //nolint:attgo_raw_string
)
