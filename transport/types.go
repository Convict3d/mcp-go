// Package transport provides MCP transport layer declarations
package transport

import "context"

// Transport interface defines the transport layer
type Transport interface {
	Call(ctx context.Context, result interface{}, method string, params ...interface{}) error
	CallRaw(ctx context.Context, method string, params interface{}) (map[string]interface{}, error)
	GetSessionID() string
	Close() error
}
