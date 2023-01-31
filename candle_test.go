package gominitrader

import (
	"fmt"
	"testing"
)

func TestMarshalCapitalPrices(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()
	capitalPricesResponse, _ := capClient.GetHistoricalPrices("USDMXN", MINUTE_30)

	var candles Candles
	err := candles.MarshalCapitalPrices(capitalPricesResponse.Prices)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	if len(candles) == 0 {
		t.Errorf("Candles Not Being Pulled or Marshalled Properly")
	}

	t.Logf("Marshalled Candles: %v+\n", candles)
}
