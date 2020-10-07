package config

import "context"

// Fetcher defines an interface that can be used for fetching data.
type Fetcher interface {
	Get(context.Context, string) ([]byte, error)
}
