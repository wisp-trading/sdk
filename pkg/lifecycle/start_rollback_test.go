package lifecycle

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	mockLifecycle "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/lifecycle"
	mockRegistry "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

func newControllerUnderTest(
	t *testing.T,
	domains []lifecycleTypes.DomainLifecycle,
	orch lifecycleTypes.Orchestrator,
) *controller {
	t.Helper()

	reg := mockRegistry.NewConnectorRegistry(t)
	reg.EXPECT().
		Filter(mock.Anything).
		Return([]connector.Connector{nil}).
		Maybe()

	return &controller{
		domains:           domains,
		connectorRegistry: reg,
		orchestrator:      orch,
		logger:            logging.NewNoOpLogger(),
		state:             lifecycleTypes.StateCreated,
		readyChan:         make(chan struct{}),
		shutdownCh:        make(chan struct{}),
	}
}

func TestStartRollsBackPriorDomainsOnFailure(t *testing.T) {
	d1 := mockLifecycle.NewDomainLifecycle(t)
	d2 := mockLifecycle.NewDomainLifecycle(t)
	orch := mockLifecycle.NewOrchestrator(t)

	d1.EXPECT().Name().Return("spot").Maybe()
	d2.EXPECT().Name().Return("perp").Maybe()
	d1.EXPECT().Start(mock.Anything, mock.Anything).Return(nil).Once()
	d2.EXPECT().Start(mock.Anything, mock.Anything).Return(errors.New("perp boom")).Once()
	d1.EXPECT().Stop().Return(nil).Once()

	ctrl := newControllerUnderTest(t, []lifecycleTypes.DomainLifecycle{d1, d2}, orch)

	err := ctrl.Start(context.Background(), strategy.StrategyName("test"), &config.StartupConfig{})
	if err == nil {
		t.Fatal("expected start error")
	}
	if ctrl.State() != lifecycleTypes.StateCreated {
		t.Fatalf("state=%v want Created", ctrl.State())
	}
}

func TestStartRollsBackDomainsWhenOrchestratorFails(t *testing.T) {
	d1 := mockLifecycle.NewDomainLifecycle(t)
	orch := mockLifecycle.NewOrchestrator(t)

	d1.EXPECT().Name().Return("spot").Maybe()
	d1.EXPECT().Start(mock.Anything, mock.Anything).Return(nil).Once()
	d1.EXPECT().Stop().Return(nil).Once()
	orch.EXPECT().Start(mock.Anything).Return(errors.New("orch boom")).Once()

	ctrl := newControllerUnderTest(t, []lifecycleTypes.DomainLifecycle{d1}, orch)

	err := ctrl.Start(context.Background(), strategy.StrategyName("test"), &config.StartupConfig{})
	if err == nil {
		t.Fatal("expected start error")
	}
	if ctrl.State() != lifecycleTypes.StateCreated {
		t.Fatalf("state=%v want Created", ctrl.State())
	}
}

func TestTriggerShutdownClosesShutdownRequested(t *testing.T) {
	ctrl := newControllerUnderTest(t, nil, mockLifecycle.NewOrchestrator(t))

	select {
	case <-ctrl.ShutdownRequested():
		t.Fatal("channel should not be closed yet")
	default:
	}

	ctrl.triggerShutdown()

	select {
	case <-ctrl.ShutdownRequested():
		// ok
	default:
		t.Fatal("expected ShutdownRequested closed after triggerShutdown")
	}

	// Idempotent
	ctrl.triggerShutdown()
}
