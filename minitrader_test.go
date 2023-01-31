package forexbot

import (
	"math"
	"testing"
	"time"
)

func _TestNewPool() *MinitraderPool {
	capClient, _ := _TestCapitalClient()
	minitrader := NewMinitrader("USDMXN", 50, 5, MINUTE_15, GPTStrategy)

	minitraderPool, _ := NewMinitraderPool(
		capClient,
		minitrader,
	)

	return minitraderPool
}

func TestMinitraderPoolUpdateCandlesData(t *testing.T) {
	pool := _TestNewPool()
	minitrader := pool.Minitraders[0]

	go pool.AuthenticateSession(time.Minute)
	go pool.UpdateMinitradersData(time.Second)

	select {
	case candles := <-minitrader.candlesChannel:
		t.Logf("Marshalled Candles: %v+\n", candles)
	case <-time.After(time.Second * 10):
		t.Error("Goroutine took too long to complete")
	}
}

func TestStrategy(t *testing.T) {
	pool := _TestNewPool()
	minitrader := pool.Minitraders[0]

	go pool.AuthenticateSession(time.Minute)
	go pool.UpdateMinitradersData(time.Second)

	select {
	case candles := <-minitrader.candlesChannel:
		signal, price := GPTStrategy(candles)
		t.Logf("Price: %s  Signal: %f", signal, price)
	case <-time.After(time.Second * 10):
		t.Error("Goroutine took too long to complete")
	}
}

func TestMinitradersAmount(t *testing.T) {
	toleranceError := 1e-5

	t.Run("2 Minitraders", func(t *testing.T) {
		capClient, _ := _TestCapitalClient()

		minitrader1 := NewMinitrader("USDMXN", 50, 5, MINUTE_15, GPTStrategy)
		minitrader2 := NewMinitrader("USDCAD", 50, 5, MINUTE_15, GPTStrategy)
		minitraderPool, _ := NewMinitraderPool(
			capClient,
			minitrader1,
			minitrader2,
		)

		tests := []struct {
			minitraderStatus []MinitraderStatus
			expectedResult   []float64
		}{
			{[]MinitraderStatus{RUNNING, RUNNING}, []float64{50, 50}},
			{[]MinitraderStatus{RUNNING, HOLDING}, []float64{100, 0}},
			{[]MinitraderStatus{HOLDING, RUNNING}, []float64{0, 100}},
			{[]MinitraderStatus{HOLDING, HOLDING}, []float64{0, 0}},
		}

		for i, test := range tests {
			minitrader1.Status, minitrader2.Status = test.minitraderStatus[0], test.minitraderStatus[1]
			minitraderPool.updateMinitradersVolatileValues(0)

			if math.Abs(minitrader1.volatileInvestmentPercentage-test.expectedResult[0]) > toleranceError {
				t.Errorf("Test case %d: minitrader1 volatile investment percentage not updated correctly, expected %f, got %f", i, test.expectedResult[0], minitrader1.volatileInvestmentPercentage)
			}
			if math.Abs(minitrader2.volatileInvestmentPercentage-test.expectedResult[1]) > toleranceError {
				t.Errorf("Test case %d: minitrader2 volatile investment percentage not updated correctly, expected %f, got %f", i, test.expectedResult[1], minitrader2.volatileInvestmentPercentage)
			}
		}
	})

	t.Run("3 Minitraders", func(t *testing.T) {
		capClient, _ := _TestCapitalClient()
		minitrader1 := NewMinitrader("USDMXN", 10, 3, MINUTE_15, GPTStrategy)
		minitrader2 := NewMinitrader("USDCAD", 30, 5, MINUTE_30, GPTStrategy)
		minitrader3 := NewMinitrader("USDJPN", 60, 5, MINUTE_5, GPTStrategy)
		minitraderPool, _ := NewMinitraderPool(capClient, minitrader1, minitrader2, minitrader3)

		tests := []struct {
			minitraderStatus []MinitraderStatus
			expectedResult   []float64
		}{
			{[]MinitraderStatus{RUNNING, RUNNING, HOLDING}, []float64{25.0, 75.0, 0}},
			{[]MinitraderStatus{RUNNING, HOLDING, HOLDING}, []float64{100.0, 0, 0}},
			{[]MinitraderStatus{RUNNING, RUNNING, RUNNING}, []float64{10.0, 30.0, 60.0}},
			{[]MinitraderStatus{HOLDING, RUNNING, RUNNING}, []float64{0, 33.333333, 66.666666}},
		}

		for i, test := range tests {
			minitrader1.Status, minitrader2.Status, minitrader3.Status = test.minitraderStatus[0], test.minitraderStatus[1], test.minitraderStatus[2]
			minitraderPool.updateMinitradersVolatileValues(0)

			if math.Abs(minitrader1.volatileInvestmentPercentage-test.expectedResult[0]) > toleranceError {
				t.Errorf("Test case %d: minitrader1 volatile investment percentage not updated correctly, expected %f, got %f", i, test.expectedResult[0], minitrader1.volatileInvestmentPercentage)
			}
			if math.Abs(minitrader2.volatileInvestmentPercentage-test.expectedResult[1]) > toleranceError {
				t.Errorf("Test case %d: minitrader2 volatile investment percentage not updated correctly, expected %f, got %f", i, test.expectedResult[1], minitrader2.volatileInvestmentPercentage)
			}
			if math.Abs(minitrader3.volatileInvestmentPercentage-test.expectedResult[2]) > toleranceError {
				t.Errorf("Test case %d: minitrader3 volatile investment percentage not updated correctly, expected %f, got %f", i, test.expectedResult[2], minitrader3.volatileInvestmentPercentage)
			}
		}
	})
}
