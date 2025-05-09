// Package mongodb backoff
package mongodb

import (
	"time"

	"github.com/jpillora/backoff"
)

func defaultBackOff() *backoff.Backoff {
	return createBackOff(100*time.Millisecond, 2*time.Minute, 2, true)
}

func createBackOff(low time.Duration, high time.Duration, factor float64, jitter bool) *backoff.Backoff {
	return &backoff.Backoff{
		Min:    low,
		Max:    high,
		Factor: factor,
		Jitter: jitter,
	}
}
