package gominitrader

import (
	"time"
)

type Candles []Candle

type Candle struct {
	Volume    int64
	Timestamp int64
	Open      BidAskPrice
	High      BidAskPrice
	Low       BidAskPrice
	Close     BidAskPrice
}

type BidAskPrice struct {
	Bid float64
	Ask float64
}

func (candles *Candles) MarshalCapitalPrices(capitalPrices []CapitalPrice) error {
	for index, capitalPrice := range capitalPrices {
		// pushing empty candle
		*candles = append(*candles, Candle{})

		// marshalling volume and timestamp fields
		time, err := time.Parse("2006-01-02T15:04:05", capitalPrice.SnapshotTimeUTC)
		if err != nil {
			return err
		}
		timestamp := time.Unix()
		(*candles)[index].Timestamp = timestamp
		(*candles)[index].Volume = int64(capitalPrice.LastTradedVolume)

		// mapping prices
		(*candles)[index].Open.Bid = capitalPrice.OpenPrice.Bid
		(*candles)[index].Open.Ask = capitalPrice.OpenPrice.Ask

		(*candles)[index].High.Bid = capitalPrice.HighPrice.Bid
		(*candles)[index].High.Ask = capitalPrice.HighPrice.Ask

		(*candles)[index].Low.Bid = capitalPrice.LowPrice.Bid
		(*candles)[index].Low.Ask = capitalPrice.LowPrice.Ask

		(*candles)[index].Close.Bid = capitalPrice.ClosePrice.Bid
		(*candles)[index].Close.Ask = capitalPrice.ClosePrice.Ask
	}
	return nil
}
