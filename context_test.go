package testext_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/ninedraft/testext"
)

func TestContext(test *testing.T) {
	test.Parallel()
	var _ testext.T = test
	var gadget = &testGadget{
		deadline: time.Now().Add(time.Hour),
	}
	var ctx = testext.Context(gadget)
	if ctx.Err() != nil {
		test.Fatal("test context already canceled: ", ctx.Err())
	}
	gadget.Fail()
	time.Sleep(time.Second)
	if ctx.Err() == nil {
		test.Fatal("test context must be canceled")
	}
}

type testGadget struct {
	deadline time.Time
	failed   int32
	toClean  []func()
}

func (t *testGadget) Helper() {}

func (t *testGadget) Fail() {
	atomic.StoreInt32(&t.failed, 1)
}

func (t *testGadget) Failed() bool {
	return atomic.LoadInt32(&t.failed) != 0
}

func (t *testGadget) Deadline() (time.Time, bool) {
	return t.deadline, !t.deadline.IsZero()
}

func (t *testGadget) Cleanup(fn func()) {
	t.toClean = append(t.toClean, fn)
}

func (t *testGadget) Finish() {
	for _, fn := range t.toClean {
		fn()
	}
}
