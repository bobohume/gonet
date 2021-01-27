package base

import "sync"

type(
	Locker struct {
		sync.Locker
	}
)
