package testext

import (
	"context"
	"time"
)

// DefaultTimeout is the default time period, after which test expected to fail.
const DefaultTimeout = 30 * time.Second

// Context creates a new context for test, which will be canceled
// after t.Deadline() or DefaultTimeout.
// Context will be eventually canceled if the test fails.
func Context(t T) context.Context {
	t.Helper()

	var deadline, deadlineOK = t.Deadline()
	if !deadlineOK {
		deadline = time.Now().Add(DefaultTimeout)
	}
	var until = time.Until(deadline)
	deadline = deadline.Add(-until / 10)
	var ctx, cancel = context.WithDeadline(context.Background(), deadline)
	t.Cleanup(cancel)
	go func() {
		defer cancel()
		var dt = until / 1000
		var timer = time.NewTimer(dt)
		defer timer.Stop()
		for !t.Failed() {
			select {
			case <-timer.C:
				timer.Reset(dt)
				continue
			case <-ctx.Done():
				return
			}
		}
	}()
	return ctx
}
