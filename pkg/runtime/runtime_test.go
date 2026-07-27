package runtime

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	mockLifecycle "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

func TestStopUsesFreshTimeoutContext(t *testing.T) {
	ctrl := mockLifecycle.NewController(t)

	var sawCanceled bool
	ctrl.EXPECT().
		Stop(mock.Anything).
		Run(func(ctx context.Context) {
			sawCanceled = ctx.Err() != nil
			_, _ = ctx.Deadline()
		}).
		Return(nil).
		Once()

	r := &rt{
		controller: ctrl,
		logger:     logging.NewNoOpLogger(),
	}
	// Simulate the old bug: root context already canceled before Stop.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r.ctx = ctx
	r.cancel = func() {}

	if err := r.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if sawCanceled {
		t.Fatal("controller.Stop received an already-canceled context")
	}
}

func TestWaitStopsOnRemoteShutdown(t *testing.T) {
	ctrl := mockLifecycle.NewController(t)
	shutdownCh := make(chan struct{})

	ctrl.EXPECT().ShutdownRequested().Return((<-chan struct{})(shutdownCh))
	ctrl.EXPECT().Stop(mock.Anything).Return(nil).Once()

	r := &rt{
		controller: ctrl,
		logger:     logging.NewNoOpLogger(),
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- r.Wait() }()

	time.Sleep(50 * time.Millisecond)
	close(shutdownCh)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after remote shutdown")
	}
}

func TestWaitRequiresStart(t *testing.T) {
	ctrl := mockLifecycle.NewController(t)
	r := &rt{controller: ctrl, logger: logging.NewNoOpLogger()}
	if err := r.Wait(); err == nil {
		t.Fatal("expected error when Wait called before Start")
	}
}
