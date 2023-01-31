package gominitrader

type Strategy func(candles Candles) (Signal, float64)

type Signal string

const (
	BUY  Signal = "BUY"
	SELL Signal = "SELL"
	NONE Signal = "NONE"
)
