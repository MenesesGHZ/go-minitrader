package gominitrader

import (
	"testing"
)

func TestMarshalCapitalPrices(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()
	capitalPricesResponse, _ := capClient.GetHistoricalPrices("USDMXN", MINUTE_15, 250)

	var candles Candles
	err := candles.MarshalCapitalPrices(capitalPricesResponse.Prices)
	if err != nil {
		t.Error(err)
	}
	if len(candles) == 0 {
		t.Error("Candles Not Being Pulled or Marshalled Properly")
	}
	if len(candles) != 250 {
		t.Errorf("Missing Candles To Pull. Current Number of Candles: %d", len(candles))
	}

	//t.Logf("Marshalled Candles: %v+\n", candles)
}
