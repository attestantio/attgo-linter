// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");

package interfacecheck

// Reader is an interface for reading.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is an interface for writing.
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Closer is an interface for closing.
type Closer interface {
	Close() error
}

// EmptyInterface has no methods - should be ignored.
type EmptyInterface interface{}

// --- Structs with interface checks (Good) ---

// GoodReader has a compliance check.
type GoodReader struct{}

var _ Reader = (*GoodReader)(nil)

func (r *GoodReader) Read(p []byte) (int, error) {
	return 0, nil
}

// --- Structs without interface checks (Bad) ---

// BadWriter implements Writer but has no compliance check.
type BadWriter struct{} // want `struct "BadWriter" implements interface "Writer"; consider adding: var _ Writer = \(\*BadWriter\)\(nil\)`

func (w *BadWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// BadCloser implements Closer but has no compliance check.
type BadCloser struct{} // want `struct "BadCloser" implements interface "Closer"; consider adding: var _ Closer = \(\*BadCloser\)\(nil\)`

func (c *BadCloser) Close() error {
	return nil
}

// --- Struct that doesn't implement any interface ---

// Helper doesn't implement any interfaces - no warning.
type Helper struct{}

func (h *Helper) DoSomething() {}
