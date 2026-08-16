package graph

import "time"

func orderCreatedAt(t *time.Time) time.Time {
	if t == nil {
		return time.Now()
	}
	return *t
}