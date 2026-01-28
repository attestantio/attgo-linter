// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");

package funcopts

import "context"

// UserService is a service type.
type UserService struct {
	db      interface{}
	cache   interface{}
	logger  interface{}
	timeout int
}

// Bad: too many parameters.
func NewUserService(db, cache, logger, metrics interface{}) *UserService { // want `constructor "NewUserService" has many parameters; consider using functional options pattern`
	return &UserService{}
}

// OrderHandler is a handler type.
type OrderHandler struct{}

// Bad: many non-context parameters.
func NewOrderHandler(ctx context.Context, db, cache, logger, validator interface{}) *OrderHandler { // want `constructor "NewOrderHandler" has many parameters; consider using functional options pattern`
	_ = ctx

	return &OrderHandler{}
}

// PaymentService uses functional options - good.
type PaymentService struct{}

// Option is a functional option.
type Option func(*PaymentService)

// Good: already uses functional options.
func NewPaymentService(opts ...Option) *PaymentService {
	return &PaymentService{}
}

// NotificationManager has few parameters - ok.
type NotificationManager struct{}

// Good: only 2 parameters.
func NewNotificationManager(db, logger interface{}) *NotificationManager {
	return &NotificationManager{}
}

// Helper is not a service type.
type Helper struct{}

// Good: not a service type, so no warning.
func NewHelper(a, b, c, d interface{}) *Helper {
	return &Helper{}
}

// Non-constructor function - not checked.
func ProcessData(a, b, c, d, e interface{}) {}
