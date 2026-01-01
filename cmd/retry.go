package cmd

import (
	"context"
	"time"
)

func withRetry(
	ctx context.Context,
	retries int,
	initialBackoff time.Duration,
	fn func() error,
) error {

	var err error
	delay := initialBackoff

	for attempt := 0; attempt <= retries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = fn()
		if err == nil {
			return nil
		}

		if attempt == retries {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			delay *= 2
		}
	}

	return err
}
