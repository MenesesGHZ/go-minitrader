package forexbot

import (
	"time"

	mapset "github.com/deckarep/golang-set"
)

type MinitraderPool struct {
	Minitraders   []Minitrader
	CapitalClient *CapitalClientAPI
}

func (pool *MinitraderPool) Start() {
	for _, minitrader := range pool.Minitraders {
		go minitrader.Start()
	}
	go pool.UpdateCandlesData(time.Second)
	go pool.UpdateMarketStatus(time.Minute)
	go pool.RenovateSession(time.Minute * 9)
}

func (pool *MinitraderPool) UpdateMarketStatus(sleepTime time.Duration) {
	for {
		// creates a tempora epic minitraders map, since multiple minitraders can be using same Epic
		epicMinitradesMap := make(map[string][]*Minitrader)

		// creates a set of epics, to then convert it to a slice of strings.
		epicsSet := mapset.NewSet()
		for _, minitrader := range pool.Minitraders {
			epicsSet.Add(minitrader.Epic)
			epicMinitradesMap[minitrader.Epic] = append(epicMinitradesMap[minitrader.Epic], &minitrader)
		}
		interfaceSlice := epicsSet.ToSlice()
		epics := make([]string, len(interfaceSlice))
		for i, v := range interfaceSlice {
			epics[i] = v.(string)
		}

		// fetch market details
		marketsDetailsResponse, err := pool.CapitalClient.GetMarketsDetails(epics)
		if err != nil {
			// send trough channel a nessage to create a new session
			//marketStatusUpdated <- false
		}
		for _, detail := range marketsDetailsResponse.MarketDetails {
			marketStatus := MinitraderMarketStatus(detail.Snapshot.MarketStatus)
			for _, minitrader := range epicMinitradesMap[detail.Instrument.Epic] {
				minitrader.MarketStatus = marketStatus
			}
		}
		// marketStatusUpdated <- true

		time.Sleep(sleepTime)
	}
}

func (pool *MinitraderPool) UpdateCandlesData(sleepTime time.Duration) {
	for {
		// build a map for avoiding requesting same data while getting historical prices
		// giving a key, the minitrader list for that key will contain minitraders
		// with the same epic and timeframe
		epicTimeframMinitraderMap := make(map[string][]*Minitrader)
		for _, minitrader := range pool.Minitraders {
			key := minitrader.Epic + string(minitrader.Timeframe)
			_, ok := epicTimeframMinitraderMap[key]
			if !ok {
				epicTimeframMinitraderMap[key] = make([]*Minitrader, 0)
			}
			epicTimeframMinitraderMap[key] = append(epicTimeframMinitraderMap[key], &minitrader)
		}

		for _, minitraders := range epicTimeframMinitraderMap {
			epic, timeframe := minitraders[0].Epic, minitraders[0].Timeframe
			pricesResponse, err := pool.CapitalClient.GetHistoricalPrices(epic, timeframe)

			var candles Candles
			candles.MarshalCapitalPrices(pricesResponse.Prices)

			for _, minitrader := range minitraders {
				if err != nil {
					// set the ministatus
					minitrader.Status = ERROR_ON_UPDATE_CANDLES_DATA
					// <-IsActive
					continue
				}
				minitrader.candlesChannel <- &candles
			}
		}

		time.Sleep(sleepTime)
	}
}

func (pool *MinitraderPool) RenovateSession(sleepTime time.Duration) {
	for {
		pool.CapitalClient.CreateNewSession()
		time.Sleep(sleepTime)
	}
}
