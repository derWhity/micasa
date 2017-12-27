// Package fhem contains a service interface and -implementation that wraps the (API-)communication with a FHEM
// instance via HTTP. It provides high-level methods to query and control devices registered in FHEM
package fhem

import (
	"context"
)

// The Service interface represents a FHEM service that can be used to interact with a FHEM instance and the devices
// registered inside of it
type Service interface {
	Discover(ctx context.Context) ([]Device, error)
}
