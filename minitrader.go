package forexbot

type MinitraderStatus string

const (
	HOLDING           MinitraderStatus = "HOLDING"
	SELL_ORDER_ACTIVE MinitraderStatus = "SELL_ORDER_ACTIVE"
	BUY_ORDER_ACTIVE  MinitraderStatus = "BUY_ORDER_ACTIVE"
)

type MinitraderMarketStatus string

const (
	TRADEABLE MinitraderMarketStatus = "TRADEABLE"
	CLOSED    MinitraderMarketStatus = "CLOSED"
)

type Minitrader struct {
	ID                   int
	Epic                 string
	InvestmentPercentage float32
	Timeframe            Timeframe
	Status               MinitraderStatus
	MarketStatus         MinitraderMarketStatus
	//Strateg
}
