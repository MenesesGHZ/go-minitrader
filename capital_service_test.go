package gominitrader

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func _TestCapitalClient() (client *CapitalClientAPI, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return client, err
	}
	capitalEmail := os.Getenv("CAPITAL_EMAIL")
	capitalApiKey := os.Getenv("CAPITAL_API_KEY")
	capitalApiKeyPassword := os.Getenv("CAPITAL_API_KEY_PASSWORD")
	capClient, err := NewCapitalClient(capitalEmail, capitalApiKey, capitalApiKeyPassword, true)
	return capClient, err
}

func TestNewCapitalClient(t *testing.T) {
	_, err := _TestCapitalClient()
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetEncriptionKey(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	encription, err := capClient.GetEncriptionKey()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(encription.EncryptionKey) == 0 {
		t.Errorf("Key is Empty")
	}
	t.Logf("Gathered Key: %s", encription.EncryptionKey)
}

func TestGetEncryptedPassword(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	encription, _ := capClient.GetEncriptionKey()
	encryptedPassword, err := capClient.GetEncryptedPassword(encription)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("Gathered Key: %s", encryptedPassword)
}

func TestCreateNewSessionAccount(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	newSessionResponse, headerTokens, err := capClient.CreateNewSession()
	if err != nil {
		t.Errorf("%v", err)
	}
	if newSessionResponse.ClientId == "" {
		t.Errorf("ClientId should be populated. Contact `gerardo.meneses.hz@gmail.com`")
	}
	t.Logf("SessionResponse: %+v", newSessionResponse)
	t.Logf("Token Headers: %+v", headerTokens)
}

func TestListWatchList(t *testing.T) { // TODO FIX, probably they change the json key
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	watchListResponse, err := capClient.GetWatchLists()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("WatchLists: %+v", watchListResponse)
}

func TestGetAllAccounts(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	accounts, err := capClient.GetAllAccounts()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("Accounts: %+v", accounts)
}

func TestGetMarketsDetails(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	marketsDetails, err := capClient.GetMarketsDetails([]string{"USDMXN", "EURUSD"})
	if err != nil {
		t.Errorf("%v", err)
	}
	if marketsDetails.MarketDetails[0].Instrument.Epic == "" {
		t.Errorf("Something is wrong with the response. Probably")
	}
	t.Logf("Market Details: %+v", marketsDetails)
}

func TestGetPrices(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	pricesResponse, err := capClient.GetHistoricalPrices("USDMXN", MINUTE_30, 250)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(pricesResponse.Prices) == 0 {
		t.Errorf("No Data Parsed. Something is wrong with the response.")
	}
	t.Logf("Market Details: %+v", pricesResponse)
}

func TestGetPositions(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	positionsResponse, err := capClient.GetPositions()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("Market Details: %+v", positionsResponse)
}

func TestCreateWorkingOrder(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	workingOrderResponse, err := capClient.CreateWorkingOrder("BTCUSD", BUY, LIMIT, 24000.0, 10)
	if err != nil {
		t.Errorf("%v", err)

	}
	t.Logf("Market Details: %+v", workingOrderResponse)
}

func TestGetAllWorkingOrders(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()
	capClient.CreateWorkingOrder("USDMXN", BUY, LIMIT, 19.20, 1000)

	workingOrdersResponse, err := capClient.GetAllWorkingOrders()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("Market Details: %+v", workingOrdersResponse)
}

func TestGetPositionOrderConfirmation(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	workingOrderResponse, _ := capClient.CreateWorkingOrder("BTCUSD", BUY, LIMIT, 24000.0, 10)
	positionOrderConfirmation, err := capClient.GetPositionOrderConfirmation(workingOrderResponse.DealReference)
	if err != nil {
		t.Errorf("%v", err)

	}
	t.Logf("Market Details: %+v", positionOrderConfirmation)
}

func TestGetPreferredAccount(t *testing.T) {
	capClient, _ := _TestCapitalClient()
	capClient.CreateNewSession()

	accountResponse, err := capClient.GetPreferredAccount()
	if err != nil {
		t.Errorf("%v", err)
	}
	if accountResponse.AccountID == "" {
		t.Error("No Data Parsed. Something is wrong with the response.")
	}
	t.Logf("Prefered Account: %+v", accountResponse)
}
