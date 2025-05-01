package utils

import (
	"time"
)

func WithRetry(retryFunc func() (bool, error)) error {
	var retryDelays = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	shouldRetry, err := retryFunc()
	if err == nil {
		return nil
	}

	if !shouldRetry {
		return err
	}

	for _, delay := range retryDelays {
		time.Sleep(delay)
		shouldRetry, err = retryFunc()
		if err == nil {
			return nil
		}

		if !shouldRetry {
			return err
		}

	}

	return err
}
