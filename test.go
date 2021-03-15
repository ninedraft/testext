package testext

import (
	"time"
)

// T is a test gadget interface.
type T interface {
	Helper()
	Failed() bool
	Deadline() (time.Time, bool)
	Cleanup(fn func())
}
