
package constants

import "time"

const (
	DefaultPort = ":8080"
	MaxBodySize = 1_048_578
)

const (
	DefaultDBMaxOpenConns = 50
	DefaultDBMaxIdleConns = 10
	DefaultDBMaxLifetime  = "15m"
	QueryTimeoutDuration  = 3 * time.Second
)

const (
	DefaultMailExpiration = 10 * time.Minute
	MinPasswordLength     = 6
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 50
	DefaultPage     = 1
)
