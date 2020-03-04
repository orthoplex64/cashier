package basecache

import (
	"time"
)

type janitor struct {
	Interval time.Duration
	Stop     chan struct{}
}

func (j *janitor) Run(c TimeAwareCache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.Stop:
			ticker.Stop()
			return
		}
	}
}
