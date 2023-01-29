package forexbot

import "fmt"

type Minitrader struct {
	Epic                 string
	InvestmentPercentage float64
	Timeframe            Timeframe
	Status               MinitraderStatus
	MarketStatus         MinitraderMarketStatus
	Strategy             Strategy
	candlesChannel       chan Candles // TODO: Implement "Pipeline" Pattern To Handle Larger Data Efficiently
}

type MinitraderStatus string

const (
	// healthy statuses
	INITIALIZING      MinitraderStatus = "INITIALIZING"
	RUNNING           MinitraderStatus = "RUNNING"
	HOLDING           MinitraderStatus = "HOLDING"
	SELL_ORDER_ACTIVE MinitraderStatus = "SELL_ORDER_ACTIVE"
	BUY_ORDER_ACTIVE  MinitraderStatus = "BUY_ORDER_ACTIVE"

	// error statuses
	ERROR_ON_UPDATE_CANDLES_DATA MinitraderStatus = "ERROR_ON_UPDATE_CANDLES_DATA"
)

type MinitraderMarketStatus string

const (
	TRADEABLE MinitraderMarketStatus = "TRADEABLE"
	CLOSED    MinitraderMarketStatus = "CLOSED"
)

type Timeframe string

const (
	MINUTE    Timeframe = "MINUTE"
	MINUTE_5  Timeframe = "MINUTE_5"
	MINUTE_15 Timeframe = "MINUTE_15"
	MINUTE_30 Timeframe = "MINUTE_30"
	HOUR      Timeframe = "HOUR"
	HOUR_4    Timeframe = "HOUR_4"
	DAY       Timeframe = "DAY"
	WEEK      Timeframe = "WEEK"
)

func NewMinitrader(epic string, investmentPercentage float64, timeframe Timeframe, strategy Strategy) *Minitrader {
	return &Minitrader{
		Epic:                 epic,
		InvestmentPercentage: investmentPercentage,
		Timeframe:            timeframe,
		Strategy:             strategy,
		Status:               INITIALIZING,
		candlesChannel:       make(chan Candles),
	}
}

func (minitrader *Minitrader) Start() {
	select {
	case candles := <-minitrader.candlesChannel:
		signal, price := minitrader.Strategy(candles)
		fmt.Println(signal, price)
	default:
		fmt.Println("hello there")
	}
}
