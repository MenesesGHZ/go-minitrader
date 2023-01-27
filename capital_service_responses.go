package forexbot

type WatchListsResponse struct {
	Epics []string `json: "epics"`
	Name  string   `json: "name"`
}

type EncriptionResponse struct {
	EncryptionKey string `json: "encryptionKey"`
	TimeStamp     int    `json: "timeStamp"`
}

type NewSessionResponse struct {
	AccountType           string      `json:"accountType"`
	AccountInfo           AccountInfo `json:"accountInfo"`
	CurrencyIsoCode       string      `json:"currencyIsoCode"`
	CurrencySymbol        string      `json:"currencySymbol"`
	CurrentAccountId      string      `json:"currentAccountId"`
	StreamingHost         string      `json:"streamingHost"`
	Accounts              []Accounts  `json:"accounts"`
	ClientId              string      `json:"clientId"`
	TimezoneOffset        int         `json:"timezoneOffset"`
	HasActiveDemoAccounts bool        `json:"hasActiveDemoAccounts"`
	HasActiveLiveAccounts bool        `json:"hasActiveLiveAccounts"`
	TrailingStopsEnabled  bool        `json:"trailingStopsEnabled"`
}

type AccountInfo struct {
	Balance    float64 `json:"balance"`
	Deposit    float64 `json:"deposit"`
	ProfitLoss float64 `json:"profitLoss"`
	Available  float64 `json:"available"`
}

type Accounts struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
	Preferred   bool   `json:"preferred"`
	AccountType string `json:"accountType"`
}
