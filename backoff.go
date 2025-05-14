package fsm

import "time"

type Backoff func(attempt int) time.Duration

func LinearBackoff(increment, max time.Duration) Backoff {
	return func(attempt int) time.Duration {
		return min(increment*time.Duration(attempt), max)
	}
}

func ExponentialBackoff(base, max time.Duration) Backoff {
	return func(attempt int) time.Duration {
		return min(base*time.Duration(1<<attempt), max)
	}
}
