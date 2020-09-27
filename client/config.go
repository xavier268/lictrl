package client

import "time"

// Configuration parameters to construct a new Client.
type Configuration struct {
	License   string // license string or number
	ServerURL string // Server entry point. Set to empty string to disable all checks.

	OfflineLimit time.Duration // how long do we accept to run without a valid check online ?
}
