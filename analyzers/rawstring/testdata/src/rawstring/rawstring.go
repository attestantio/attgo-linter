// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");

package rawstring

// Bad: multiple escapes (4 backslash pairs).
var path1 = "C:\\Users\\name\\Documents\\file.txt" // want `string has 4 escape sequences; consider using a raw string`

// Good: already a raw string.
var query2 = `vouch_relay_execution_config_total{result="succeeded"}`

// Good: only two escapes - below threshold (3).
var simple = "hello \"world\""

// Good: also only two escapes.
var query1 = "vouch_relay_execution_config_total{result=\"succeeded\"}"

// Good: no escapes.
var plain = "hello world"

// Good: \n serves a purpose (newline) - not counted.
var multiline = "line1\nline2"

// Good: \t serves a purpose (tab) - not counted.
var tabbed = "col1\tcol2"

// Good: contains backtick - can't use raw string.
var withBacktick = "use `backticks` here"

// Good: raw string with quotes inside.
var rawWithQuotes = `she said "hello"`

// Bad: JSON with escapes (6 escape sequences: 3 pairs of \").
var jsonStr = "{\"key\": \"value\", \"num\": 123}" // want `string has 6 escape sequences; consider using a raw string`

func examples() {
	// Local variables also checked (3 backslash escapes).
	_ = "path\\to\\file\\name" // want `string has 3 escape sequences; consider using a raw string`

	// Good: simple string.
	_ = "hello"

	// Good: raw string.
	_ = `path\to\file`
}
