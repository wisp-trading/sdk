package lifecycle

import (
	"context"
	"testing"
	"time"

	mockIngestors "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	mockHealth "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/health"
	mockLifecycle "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	mockRegistry "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
)

func TestLifecycleController_StateTransitions(t *testing.T) {
	mockMarket := mockIngestors.NewMarketDataCoordinator(t)
	mockPosition := mockIngestors.NewPositionCoordinator(t)
	mockReg := mockRegistry.NewConnectorRegistry(t)
	mockHealthStore := mockHealth.NewHealthStore(t)
	mockOrchestrator := mockLifecycle.NewOrchestrator(t)
	noopLog := logging.NewNoOpLogger()

	// Setup expectations
	mockMarket.EXPECT().StartDataCollection(context.Background()).Return(nil).Once()
	mockPosition.EXPECT().Start(context.Background()).Return(nil).Once()
	mockOrchestrator.EXPECT().Start(context.Background()).Return(nil).Once()
	mockReg.EXPECT().GetReadyConnectors().Return([]connector.Connector{}).Maybe()
	mockMarket.EXPECT().StopDataCollection().Return(nil).Once()
	mockPosition.EXPECT().Stop().Return(nil).Once()
	mockOrchestrator.EXPECT().Stop(context.Background()).Return(nil).Once()

	controller := NewController(mockMarket, mockPosition, mockReg, mockHealthStore, mockOrchestrator, noopLog)

	// Initial state should be Created
	if controller.State() != lifecycleTypes.StateCreated {
		t.Errorf("Expected initial state to be Created, got %v", controller.State())
	}

	// Should not be ready initially
	if controller.IsReady() {
		t.Error("Controller should not be ready initially")
	}

	// Start the controller
	ctx := context.Background()
	if err := controller.Start(ctx); err != nil {
		t.Fatalf("Failed to start controller: %v", err)
	}

	// State should now be Ready
	if controller.State() != lifecycleTypes.StateReady {
		t.Errorf("Expected state to be Ready after start, got %v", controller.State())
	}

	// Should be ready now
	if !controller.IsReady() {
		t.Error("Controller should be ready after starting")
	}

	// Stop the controller
	if err := controller.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop controller: %v", err)
	}

	// State should now be Stopped
	if controller.State() != lifecycleTypes.StateStopped {
		t.Errorf("Expected state to be Stopped after stop, got %v", controller.State())
	}
}

func TestLifecycleController_WaitUntilReady(t *testing.T) {
	mockMarket := mockIngestors.NewMarketDataCoordinator(t)
	mockPosition := mockIngestors.NewPositionCoordinator(t)
	mockReg := mockRegistry.NewConnectorRegistry(t)
	mockHealthStore := mockHealth.NewHealthStore(t)
	mockOrchestrator := mockLifecycle.NewOrchestrator(t)
	noopLog := logging.NewNoOpLogger()

	// Setup expectations
	mockMarket.EXPECT().StartDataCollection(context.Background()).Return(nil).Once()
	mockPosition.EXPECT().Start(context.Background()).Return(nil).Once()
	mockOrchestrator.EXPECT().Start(context.Background()).Return(nil).Once()
	mockReg.EXPECT().GetReadyConnectors().Return([]connector.Connector{}).Maybe()

	controller := NewController(mockMarket, mockPosition, mockReg, mockHealthStore, mockOrchestrator, noopLog)

	// Start in background
	ctx := context.Background()
	go func() {
		time.Sleep(100 * time.Millisecond)
		controller.Start(ctx)
	}()

	// Wait for ready
	waitCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := controller.WaitUntilReady(waitCtx); err != nil {
		t.Fatalf("Failed to wait for ready: %v", err)
	}

	if !controller.IsReady() {
		t.Error("Controller should be ready after WaitUntilReady returns")
	}
}

func TestLifecycleController_CannotStartTwice(t *testing.T) {
	mockMarket := mockIngestors.NewMarketDataCoordinator(t)
	mockPosition := mockIngestors.NewPositionCoordinator(t)
	mockReg := mockRegistry.NewConnectorRegistry(t)
	mockHealthStore := mockHealth.NewHealthStore(t)
	mockOrchestrator := mockLifecycle.NewOrchestrator(t)
	noopLog := logging.NewNoOpLogger()

	// Setup expectations
	mockMarket.EXPECT().StartDataCollection(context.Background()).Return(nil).Once()
	mockPosition.EXPECT().Start(context.Background()).Return(nil).Once()
	mockOrchestrator.EXPECT().Start(context.Background()).Return(nil).Once()
	mockReg.EXPECT().GetReadyConnectors().Return([]connector.Connector{}).Maybe()

	controller := NewController(mockMarket, mockPosition, mockReg, mockHealthStore, mockOrchestrator, noopLog)

	ctx := context.Background()

	// First start should succeed
	if err := controller.Start(ctx); err != nil {
		t.Fatalf("First start failed: %v", err)
	}

	// Second start should fail
	if err := controller.Start(ctx); err == nil {
		t.Error("Expected error when starting controller twice, got nil")
	}
}
