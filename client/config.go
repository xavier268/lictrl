package client

import "time"

// Configuration parameters to construct a new Client.
type Configuration struct {
	License   string // license string or number, includes version number, licence serial number, etc ...
	ServerURL string // Server entry point. Set to empty string to disable all checks.

	OfflineLimit time.Duration // how long do we accept to run without a valid check online ?
	AutoRepeat   time.Duration // Automatically check with the provided period. Ignore if value is zero-value (not set).
}
