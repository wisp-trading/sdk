package executor

import (
	"testing"

	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

func TestSideFromAction(t *testing.T) {
	if SideFromAction(strategy.ActionBuy) != connector.OrderSideBuy {
		t.Fatal("buy")
	}
	if SideFromAction(strategy.ActionSell) != connector.OrderSideSell {
		t.Fatal("sell")
	}
	if SideFromAction(strategy.ActionSellShort) != connector.OrderSideSell {
		t.Fatal("short")
	}
}
