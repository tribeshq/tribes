package retry

import (
	"log/slog"
	"time"
)

// Implements a simple method call retry policy.
// This retry policy will retry a failed execution for up to 'maxRetries' times
// In between retry calls it will wait a "maxDelay"
func CallFunctionWithRetryPolicy[
	R any,
	A any,
](
	fn func(A) (R, error),
	args A,
	logger *slog.Logger,
	maxRetries uint64,
	maxDelay time.Duration,
	infoLabel string,
) (R, error) {

	var lastErr error
	var lastValue R

	for i := uint64(0); i <= maxRetries; i++ {

		if i != 0 {
			logger.Info("Retry Policy: Retrying...", "delay", maxDelay)
			time.Sleep(maxDelay)
		}

		lastValue, lastErr = fn(args)
		if lastErr == nil {
			return lastValue, nil
		}
		logger.Info(
			"Retry Policy: Got error calling function",
			"label",
			infoLabel,
			"error",
			lastErr.Error())

	}
	return lastValue, lastErr

}
