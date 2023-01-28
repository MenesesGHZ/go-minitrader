package forexbot

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func testCapitalClient() (client *CapitalClientAPI, err error) {
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
	_, err := testCapitalClient()
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
}

func TestGetEncriptionKey(t *testing.T) {
	capClient, _ := testCapitalClient()
	encription, err := capClient.GetEncriptionKey()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Got Error")
	}
	if len(encription.EncryptionKey) == 0 {
		t.Errorf("Key is Empty")
	}
	t.Logf("Gathered Key: %s", encription.EncryptionKey)
}

func TestGetEncryptedPassword(t *testing.T) {
	capClient, _ := testCapitalClient()
	encription, _ := capClient.GetEncriptionKey()
	encryptedPassword, err := capClient.GetEncryptedPassword(encription)
	if err != nil {
		fmt.Println(err)
		t.Errorf("Got Error")
	}
	t.Logf("Gathered Key: %s", encryptedPassword)
}

func TestCreateNewSessionAccount(t *testing.T) {
	capClient, _ := testCapitalClient()
	newSessionResponse, headerTokens, err := capClient.CreateNewSession()
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	if newSessionResponse.ClientId == "" {
		fmt.Println("ClientId should be populated. Contact `gerardo.meneses.hz@gmail.com`")
		t.Error()
	}
	t.Logf("SessionResponse: %+v", newSessionResponse)
	t.Logf("Token Headers: %+v", headerTokens)
}

func TestListWatchList(t *testing.T) { // TODO FIX, probably they change the json key
	capClient, _ := testCapitalClient()
	capClient.CreateNewSession()

	watchListResponse, err := capClient.GetWatchLists()
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	t.Logf("WatchLists: %+v", watchListResponse)
}

func TestGetAllAccounts(t *testing.T) {
	capClient, _ := testCapitalClient()
	capClient.CreateNewSession()

	accounts, err := capClient.GetAllAccounts()
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	t.Logf("Accounts: %+v", accounts)
}

func TestGetMarkets(t *testing.T) {
	capClient, _ := testCapitalClient()
	capClient.CreateNewSession()

	marketsDetails, err := capClient.GetMarketsDetails([]string{"USDMXN", "EURUSD"})
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	if marketsDetails.MarketDetails[0].Instrument.Epic == "" {
		t.Errorf("Something is wrong with the response. Probably")
	}
	t.Logf("Market Details: %+v", marketsDetails)
}

func TestGetPrices(t *testing.T) {
	capClient, _ := testCapitalClient()
	capClient.CreateNewSession()

	pricesResponse, err := capClient.GetPrices("USDMXN", "MINUTE_30")
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	if len(pricesResponse.Prices) == 0 {
		t.Errorf("Something is wrong with the response. Probably")
	}
	t.Logf("Market Details: %+v", pricesResponse)
}

func TestGetPositions(t *testing.T) {
	capClient, _ := testCapitalClient()
	capClient.CreateNewSession()

	positionsResponse, err := capClient.GetPositions()
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	t.Logf("Market Details: %+v", positionsResponse)
}

func TestCreatePosition(t *testing.T) {
	capClient, _ := testCapitalClient()
	capClient.CreateNewSession()

	positionsResponse, err := capClient.CreatePosition("USDMXN", BUY, 19.20)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	t.Logf("Market Details: %+v", positionsResponse)
}
