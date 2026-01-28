// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");

package capitalcomment

// this is a bad comment // want `comment should start with a capital letter`

// This is a good comment

// someFunc is used to do things

// myVariable contains the data

// my_variable is set here

// nolint:errcheck

// TODO: fix this later

// See https://example.com for more info

// +build linux

// go:generate stringer -type=Foo

// ... continued from above

// 123 is the magic number

// see the documentation // want `comment should start with a capital letter`

var x = 1
