// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");

package structfieldorder

import "sync"

// GoodService has fields in the correct order.
type GoodService struct {
	// Logger
	log interface{}

	// Metrics
	metrics interface{}

	// Dependencies
	client interface{}
	db     interface{}

	// Data
	config interface{}
	name   string

	// Sync
	mu   sync.Mutex
	done chan struct{}
}

// BadService has fields out of order.
type BadService struct {
	mu      sync.Mutex  // sync should be last
	log     interface{} // want `field "log" \(logger\) should come before "mu" \(synchronization\)`
	config  interface{}
	metrics interface{} // want `field "metrics" \(metrics\) should come before "config" \(data\)`
}

// AnotherBad has sync before dependencies.
type AnotherBad struct {
	log    interface{}
	done   chan struct{} // sync
	client interface{}   // want `field "client" \(dependency\) should come before "done" \(synchronization\)`
}

// SimpleStruct with just data fields is fine.
type SimpleStruct struct {
	Name  string
	Value int
}

// MixedOrder with multiple violations.
type MixedOrder struct {
	wg      sync.WaitGroup // sync first - bad
	logger  interface{}    // want `field "logger" \(logger\) should come before "wg" \(synchronization\)`
	db      interface{}
	metrics interface{} // want `field "metrics" \(metrics\) should come before "db" \(dependency\)`
}
