package forexbot

import (
	"testing"
	"time"
)

func _TestNewPool() *MinitraderPool {
	capClient, _ := _TestCapitalClient()
	minitrader := NewMinitrader("USDMXN", 50, MINUTE_15, GPTStrategy)

	minitraderPool := NewMinitraderPool(
		capClient,
		minitrader,
	)

	return minitraderPool
}

func TestMinitraderPoolUpdateCandlesData(t *testing.T) {
	pool := _TestNewPool()
	minitrader := pool.Minitraders[0]

	go pool.AuthenticateSession(time.Minute)
	go pool.UpdateCandlesData(time.Second)

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
	go pool.UpdateCandlesData(time.Second)

	select {
	case candles := <-minitrader.candlesChannel:
		signal, price := GPTStrategy(candles)
		t.Logf("Price: %s  Signal: %f", signal, price)
	case <-time.After(time.Second * 10):
		t.Error("Goroutine took too long to complete")
	}
}
