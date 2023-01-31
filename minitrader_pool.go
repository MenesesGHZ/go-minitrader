package gominitrader

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
)

type MinitraderPool struct {
	Minitraders   []*Minitrader
	CapitalClient *CapitalClientAPI

	wg                         *sync.WaitGroup
	epics                      []string                 // slice of unique epics use on minitraders
	epicMinitraderMap          map[string][]*Minitrader // used for checking market status
	epicTimeframeMinitraderMap map[string][]*Minitrader // used for fetching historical prices
}

func NewMinitraderPool(capitalClient *CapitalClientAPI, minitraders ...*Minitrader) (*MinitraderPool, error) {
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
	// giving a key, the minitrader list for that key will contain minitraders with the same epic and timeframe
	availablePercentage := 0.0
	for _, minitrader := range minitraders {
		key := minitrader.Epic + string(minitrader.Timeframe)
		pool.epicTimeframeMinitraderMap[key] = append(pool.epicTimeframeMinitraderMap[key], minitrader)

		availablePercentage += minitrader.InvestmentPercentage
	}

	if availablePercentage != 100.0 {
		return &MinitraderPool{}, errors.New(fmt.Sprintf("Minitraders InvestmentPercentage` Sum Must Be 100.0; Current Sum: %f", availablePercentage))
	}

	return pool, nil
}

func (pool *MinitraderPool) Start() {
	for _, minitrader := range pool.Minitraders {
		minitrader.capitalClient = pool.CapitalClient
		go minitrader.Start(pool.wg)
		pool.wg.Add(1)
	}
	go pool.UpdateMinitradersData(time.Second)
	go pool.UpdateMarketStatus(time.Minute)
	go pool.AuthenticateSession(time.Minute * 9)
	go pool.Pulse()
	pool.wg.Wait()
}

func (pool *MinitraderPool) UpdateMarketStatus(sleepTime time.Duration) {
	for {
		marketsDetailsResponse, err := pool.CapitalClient.GetMarketsDetails(pool.epics)
		if _, ok := err.(*CapitalClientUnathenticated); ok {
			// sleep and retry. AuthenticateSession goroutine should handle this
			time.Sleep(sleepTime)
			continue
		} else if err != nil {
			log.Fatalf("Unexpected Error: %v", err) // TODO: Improve error handling
		}
		for _, detail := range marketsDetailsResponse.MarketDetails {
			marketStatus := MinitraderMarketStatus(detail.Snapshot.MarketStatus)
			for _, minitrader := range pool.epicMinitraderMap[detail.Instrument.Epic] {
				minitrader.MarketStatus = marketStatus
			}
		}
		time.Sleep(sleepTime)
	}
}

func (pool *MinitraderPool) UpdateMinitradersData(sleepTime time.Duration) {
	for {
		// update minitraderes amountAvailable to invest
		account, err := pool.CapitalClient.GetPreferredAccount()
		if _, ok := err.(*CapitalClientUnathenticated); ok {
			// sleep and retry. AuthenticateSession goroutine should handle this
			time.Sleep(sleepTime)
			continue
		} else if err != nil {
			log.Fatalf("Unexpected Error: %v", err) // TODO: Improve error handling
		}
		pool.updateMinitradersVolatileValues(account.Balance.Available)

		// update minitraders candles data
		for _, minitraders := range pool.epicTimeframeMinitraderMap {
			epic, timeframe := minitraders[0].Epic, minitraders[0].Timeframe
			pricesResponse, err := pool.CapitalClient.GetHistoricalPrices(epic, timeframe)
			if _, ok := err.(*CapitalClientUnathenticated); ok {
				// break loop, then sleep and retry. AuthenticateSession goroutine should handle this
				break
			} else if err != nil {
				log.Fatalf("Unexpected Error: %v", err) // TODO: Improve error handling
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
	tryCounter := 0
	for tryCounter < 3 {
		_, _, err := pool.CapitalClient.CreateNewSession()
		time.Sleep(sleepTime)
		if err != nil {
			tryCounter++
		} else {
			tryCounter = 0
		}
	}

	// stop minitrader_pool; TODO: improve logging
	for i := 0; i < len(pool.Minitraders); i++ {
		pool.wg.Done()
	}
}

func (pool *MinitraderPool) Pulse() {
	for {
		log.Print("beat.")
		time.Sleep(time.Second)
	}
}

func (pool *MinitraderPool) updateMinitradersVolatileValues(amountAvailable float64) {
	var totalPercent float64
	for _, minitrader := range pool.Minitraders {
		if minitrader.Status == NEW || minitrader.Status == RUNNING {
			totalPercent += minitrader.InvestmentPercentage
		}
	}
	for _, minitrader := range pool.Minitraders {
		if minitrader.Status != NEW && minitrader.Status != RUNNING {
			minitrader.volatileInvestmentPercentage = 0
			minitrader.volatileAmountAvailable = 0
			continue
		}
		minitrader.volatileInvestmentPercentage = minitrader.InvestmentPercentage / totalPercent * 100
		minitrader.volatileAmountAvailable = minitrader.InvestmentPercentage / 100 * amountAvailable
	}
}
