package forexbot

import (
	mapset "github.com/deckarep/golang-set"
)

type MinitraderPool struct {
	Minitraders   []Minitrader
	CapitalClient *CapitalClientAPI
}

func (pool MinitraderPool) UpdateMarketStatus() {
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
}
