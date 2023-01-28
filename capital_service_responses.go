package forexbot

import "time"

type WatchListsResponse struct {
	Epics []string `json: "epics"`
	Name  string   `json: "name"`
}

type EncriptionResponse struct {
	EncryptionKey string `json: "encryptionKey"`
	TimeStamp     int    `json: "timeStamp"`
}

type NewSessionResponse struct {
	AccountType           string `json:"accountType"`
	CurrencyIsoCode       string `json:"currencyIsoCode"`
	CurrencySymbol        string `json:"currencySymbol"`
	CurrentAccountId      string `json:"currentAccountId"`
	StreamingHost         string `json:"streamingHost"`
	ClientId              string `json:"clientId"`
	TimezoneOffset        int    `json:"timezoneOffset"`
	HasActiveDemoAccounts bool   `json:"hasActiveDemoAccounts"`
	HasActiveLiveAccounts bool   `json:"hasActiveLiveAccounts"`
	TrailingStopsEnabled  bool   `json:"trailingStopsEnabled"`
	Accounts              []struct {
		AccountId   string `json:"accountId"`
		AccountName string `json:"accountName"`
		Preferred   bool   `json:"preferred"`
		AccountType string `json:"accountType"`
	} `json:"accounts"`
	AccountInfo struct {
		Balance    float64 `json:"balance"`
		Deposit    float64 `json:"deposit"`
		ProfitLoss float64 `json:"profitLoss"`
		Available  float64 `json:"available"`
	} `json:"accountInfo"`
}

type AccountsResponse struct {
	Accounts []struct {
		AccountID   string `json:"accountId"`
		AccountName string `json:"accountName"`
		Status      string `json:"status"`
		AccountType string `json:"accountType"`
		Preferred   bool   `json:"preferred"`
		Balance     struct {
			Balance    float64 `json:"balance"`
			Deposit    float64 `json:"deposit"`
			ProfitLoss float64 `json:"profitLoss"`
			Available  float64 `json:"available"`
		} `json:"balance"`
		Currency string `json:"currency"`
	} `json:"accounts"`
}

type MarketsDetailsResponse struct {
	MarketDetails []struct {
		Instrument struct {
			Epic                     string  `json:"epic"`
			Expiry                   string  `json:"expiry"`
			Name                     string  `json:"name"`
			LotSize                  int     `json:"lotSize"`
			Type                     string  `json:"type"`
			GuaranteedStopAllowed    bool    `json:"guaranteedStopAllowed"`
			StreamingPricesAvailable bool    `json:"streamingPricesAvailable"`
			Currency                 string  `json:"currency"`
			MarginFactor             float64 `json:"marginFactor"`
			MarginFactorUnit         string  `json:"marginFactorUnit"`
			OpeningHours             string  `json:"openingHours"`
			Country                  string  `json:"country"`
		} `json:"instrument"`
		DealingRules struct {
			MinStepDistance struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"minStepDistance"`
			MinDealSize struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"minDealSize"`
			MaxDealSize struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"maxDealSize"`
			MinSizeIncrement struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"minSizeIncrement"`
			MinGuaranteedStopDistance struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"minGuaranteedStopDistance"`
			MinStopOrProfitDistance struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"minStopOrProfitDistance"`
			MaxStopOrProfitDistance struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"maxStopOrProfitDistance"`
			MarketOrderPreference   string `json:"marketOrderPreference"`
			TrailingStopsPreference string `json:"trailingStopsPreference"`
		} `json:"dealingRules"`
		Snapshot struct {
			MarketStatus        string  `json:"marketStatus"`
			UpdateTime          string  `json:"updateTime"`
			DelayTime           int     `json:"delayTime"`
			Bid                 float64 `json:"bid"`
			Offer               float64 `json:"offer"`
			DecimalPlacesFactor int     `json:"decimalPlacesFactor"`
			ScalingFactor       int     `json:"scalingFactor"`
		} `json:"snapshot"`
	} `json:"marketDetails"`
}

type PricesResponse struct {
	Prices []struct {
		SnapshotTime    string `json:"snapshotTime"`
		SnapshotTimeUTC string `json:"snapshotTimeUTC"`
		OpenPrice       struct {
			Bid float64 `json:"bid"`
			Ask float64 `json:"ask"`
		} `json:"openPrice"`
		ClosePrice struct {
			Bid float64 `json:"bid"`
			Ask float64 `json:"ask"`
		} `json:"closePrice"`
		HighPrice struct {
			Bid float64 `json:"bid"`
			Ask float64 `json:"ask"`
		} `json:"highPrice"`
		LowPrice struct {
			Bid float64 `json:"bid"`
			Ask float64 `json:"ask"`
		} `json:"lowPrice"`
		LastTradedVolume int `json:"lastTradedVolume"`
	}
}

type PositionsResponse struct {
	Positions []struct {
		Position struct {
			ContractSize   int       `json:"contractSize"`
			CreatedDate    time.Time `json:"createdDate"`
			CreatedDateUTC time.Time `json:"createdDateUTC"`
			DealID         string    `json:"dealId"`
			DealReference  string    `json:"dealReference"`
			Size           int       `json:"size"`
			Direction      string    `json:"direction"`
			Level          float64   `json:"level"`
			Currency       string    `json:"currency"`
			GuaranteedStop bool      `json:"guaranteedStop,omitempty"`
			ControlledRisk bool      `json:"controlledRisk,omitempty"`
		} `json:"position"`
		Market struct {
			InstrumentName       string    `json:"instrumentName"`
			Expiry               string    `json:"expiry"`
			MarketStatus         string    `json:"marketStatus"`
			Epic                 string    `json:"epic"`
			InstrumentType       string    `json:"instrumentType"`
			LotSize              int       `json:"lotSize"`
			High                 float64   `json:"high"`
			Low                  float64   `json:"low"`
			PercentageChange     float64   `json:"percentageChange"`
			NetChange            float64   `json:"netChange"`
			Bid                  float64   `json:"bid"`
			Offer                float64   `json:"offer"`
			UpdateTime           time.Time `json:"updateTime"`
			UpdateTimeUTC        time.Time `json:"updateTimeUTC"`
			DelayTime            int       `json:"delayTime"`
			StreamingPricesAvail bool      `json:"streamingPricesAvailable"`
			ScalingFactor        int       `json:"scalingFactor"`
		} `json:"market"`
	} `json:"positions"`
}

type CreatePositionResponse struct {
	DealReference string `json:"dealReference"`
}
