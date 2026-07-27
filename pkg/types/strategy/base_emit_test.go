package strategy

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

type stubSignal struct{ id uuid.UUID }

func (s stubSignal) GetID() uuid.UUID          { return s.id }
func (s stubSignal) GetStrategy() StrategyName { return "stub" }
func (s stubSignal) GetTimestamp() time.Time   { return time.Time{} }

func TestBaseStrategyEmitDoesNotDropWhenConsumerSlow(t *testing.T) {
	b := NewBaseStrategy(BaseStrategyConfig{Name: "emit-test"})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := b.StartWithRunner(ctx, func(ctx context.Context) { <-ctx.Done() }); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = b.Stop(context.Background()) }()

	done := make(chan struct{})
	go func() {
		for i := 0; i < signalChannelBufferSize+1; i++ {
			b.Emit(stubSignal{id: uuid.New()})
		}
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	drained := 0
	for drained < signalChannelBufferSize+1 {
		select {
		case <-b.Signals():
			drained++
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting for emits; drained=%d", drained)
		}
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Emit sender did not complete after drain")
	}
}

func TestBaseStrategyEmitRespectsCancel(t *testing.T) {
	b := NewBaseStrategy(BaseStrategyConfig{Name: "emit-cancel"})
	ctx, cancel := context.WithCancel(context.Background())

	if err := b.StartWithRunner(ctx, func(ctx context.Context) { <-ctx.Done() }); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < signalChannelBufferSize; i++ {
		b.Emit(stubSignal{id: uuid.New()})
	}

	blocked := make(chan struct{})
	go func() {
		b.Emit(stubSignal{id: uuid.New()})
		close(blocked)
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()
	_ = b.Stop(context.Background())

	select {
	case <-blocked:
	case <-time.After(2 * time.Second):
		t.Fatal("Emit did not unblock on context cancel")
	}
}
