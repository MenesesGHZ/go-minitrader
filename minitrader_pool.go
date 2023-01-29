package forexbot

import (
	"time"

	mapset "github.com/deckarep/golang-set"
)

type MinitraderPool struct {
	Minitraders   []*Minitrader
	CapitalClient *CapitalClientAPI

	epics                      []string                 // slice of unique epics use on minitraders
	epicMinitraderMap          map[string][]*Minitrader // used for checking market status
	epicTimeframeMinitraderMap map[string][]*Minitrader // used for fetching historical prices
}

func NewMinitraderPool(capitalClient *CapitalClientAPI, minitraders ...*Minitrader) *MinitraderPool {
	pool := &MinitraderPool{
		CapitalClient: capitalClient,
		Minitraders:   minitraders,

		epicMinitraderMap:          make(map[string][]*Minitrader),
		epicTimeframeMinitraderMap: make(map[string][]*Minitrader),
	}

	// creates an epic minitraders map, since multiple minitraders can be using same Epic + a slice of unique epics
	epicsSet := mapset.NewSet()
	for _, minitrader := range minitraders {
		epicsSet.Add(minitrader.Epic)
		pool.epicMinitraderMap[minitrader.Epic] = append(pool.epicMinitraderMap[minitrader.Epic], minitrader)
	}
	for _, v := range epicsSet.ToSlice() {
		pool.epics = append(pool.epics, v.(string))
	}

	// build a map for avoiding requesting same data while getting historical prices
	// giving a key, the minitrader list for that key will contain minitraders
	// with the same epic and timeframe
	for _, minitrader := range minitraders {
		key := minitrader.Epic + string(minitrader.Timeframe)
		pool.epicTimeframeMinitraderMap[key] = append(pool.epicTimeframeMinitraderMap[key], minitrader)
	}

	return pool
}

func (pool *MinitraderPool) Start() {
	for _, minitrader := range pool.Minitraders {
		go minitrader.Start()
	}
	go pool.UpdateCandlesData(time.Second)
	go pool.UpdateMarketStatus(time.Minute)
	go pool.AuthenticateSession(time.Minute * 9)
}

func (pool *MinitraderPool) UpdateMarketStatus(sleepTime time.Duration) {
	for {
		// fetch market details
		marketsDetailsResponse, err := pool.CapitalClient.GetMarketsDetails(pool.epics)
		if err != nil {
			// send trough channel a nessage to create a new session
			//marketStatusUpdated <- false
		}
		for _, detail := range marketsDetailsResponse.MarketDetails {
			marketStatus := MinitraderMarketStatus(detail.Snapshot.MarketStatus)
			for _, minitrader := range pool.epicMinitraderMap[detail.Instrument.Epic] {
				minitrader.MarketStatus = marketStatus
			}
		}
		// marketStatusUpdated <- true

		time.Sleep(sleepTime)
	}
}

func (pool *MinitraderPool) UpdateCandlesData(sleepTime time.Duration) {
	for {
		for _, minitraders := range pool.epicTimeframeMinitraderMap {
			epic, timeframe := minitraders[0].Epic, minitraders[0].Timeframe
			pricesResponse, err := pool.CapitalClient.GetHistoricalPrices(epic, timeframe)
			if _, ok := err.(*CapitalClientUnathenticated); ok {
				// break loop and retry after sleeptime ends.
				// AuthenticateSession goroutine should handle this
				break
			}

			var candles Candles
			candles.MarshalCapitalPrices(pricesResponse.Prices)

			for _, minitrader := range minitraders {
				if err != nil {
					minitrader.Status = ERROR_ON_UPDATE_CANDLES_DATA
					continue
				}
				minitrader.candlesChannel <- candles
			}
		}
		time.Sleep(sleepTime)
	}
}

func (pool *MinitraderPool) AuthenticateSession(sleepTime time.Duration) {
	for {
		pool.CapitalClient.CreateNewSession()
		time.Sleep(sleepTime)
	}
}
