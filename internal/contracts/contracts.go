package contracts

import "time"

type Bucket interface {
	Allow() bool
	LastUse() time.Time
	Period() time.Duration
	Stop()
}
