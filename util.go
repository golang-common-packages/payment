package payment

import (
	"context"
)

// SetContext set new context
func SetContext(context context.Context) {
	ctx = context
}

// GetContext return the current context
func GetContext() context.Context {
	return ctx
}
