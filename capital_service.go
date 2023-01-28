package forexbot

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const CAPITAL_DOMAIN_NAME = "https://api-capital.backend-capital.com"
const CAPITAL_DEBUG_DOMAIN_NAME = "https://demo-api-capital.backend-capital.com"

type CapitalClientAPI struct {
	CAPITAL_EMAIL            string
	CAPITAL_API_KEY          string
	CAPITAL_API_KEY_PASSWORD string
	CapitalDomainName        string
	HttpClient               *http.Client
}

func NewCapitalClient(capitalEmail string, capitalApiKey string, capitalApiKeyPassword string, debug bool) (client *CapitalClientAPI, err error) {
	if capitalEmail == "" {
		return client, errors.New("Capital Email cannot be an empty string")
	}
	if capitalApiKey == "" {
		return client, errors.New("Capital Api Key cannot be an empty string")
	}
	if capitalApiKeyPassword == "" {
		return client, errors.New("Capital Api Key Password cannot be an empty string")
	}
	capitalDomainName := CAPITAL_DOMAIN_NAME
	if debug {
		capitalDomainName = CAPITAL_DEBUG_DOMAIN_NAME
	}

	return &CapitalClientAPI{
		CAPITAL_EMAIL:            capitalEmail,
		CAPITAL_API_KEY:          capitalApiKey,
		CAPITAL_API_KEY_PASSWORD: capitalApiKeyPassword,
		CapitalDomainName:        capitalDomainName,
		HttpClient:               &http.Client{Transport: nil},
	}, nil
}

func (capClient *CapitalClientAPI) GetEncriptionKey() (EncriptionResponse, error) {
	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/session/encryptionKey", nil)
	request.Header.Add("X-CAP-API-KEY", capClient.CAPITAL_API_KEY)

	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return EncriptionResponse{}, err
	}
	defer response.Body.Close()

	var encription EncriptionResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&encription)

	return encription, nil
}

func (capClient *CapitalClientAPI) GetWatchLists() (WatchListsResponse, error) {
	if capClient.HttpClient.Transport == nil {
		return WatchListsResponse{}, errors.New("A Session is needed; Run `capClient.CreateNewSession()` to authenticate")
	}

	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/watchlists", nil)
	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return WatchListsResponse{}, err
	}
	defer response.Body.Close()

	var watchlists WatchListsResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&watchlists)

	return watchlists, nil
}

func (capClient *CapitalClientAPI) CreateNewSession() (newSessionResponse NewSessionResponse, headerTokens http.Header, err error) {
	encriptionResponse, err := capClient.GetEncriptionKey()
	if err != nil {
		return newSessionResponse, headerTokens, err
	}
	encriptedPassword, err := capClient.GetEncryptedPassword(encriptionResponse)
	if err != nil {
		return newSessionResponse, headerTokens, err
	}

	body := NewSessionBody{
		Identifier:        capClient.CAPITAL_EMAIL,
		Password:          encriptedPassword,
		EncryptedPassword: true,
	}
	jsonData, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "https://api-capital.backend-capital.com/api/v1/session", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-CAP-API-KEY", capClient.CAPITAL_API_KEY)
	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return newSessionResponse, headerTokens, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return newSessionResponse, headerTokens, errors.New(fmt.Sprintf("Unexpected [%d] Status Code Response - %s", response.StatusCode, string(body)))
	}

	// set new session resposne
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&newSessionResponse)

	// set header tokens
	headerTokens = http.Header{}
	headerTokens.Add("CST", response.Header.Get("CST"))
	headerTokens.Add("X-SECURITY-TOKEN", response.Header.Get("X-SECURITY-TOKEN"))

	// update http client transport to set new auth creds in header for new requests
	capClient.HttpClient = &http.Client{
		Transport: &AuthenticationTransport{
			RoundTripper:     http.DefaultTransport,
			X_SECURITY_TOKEN: response.Header.Get("X-SECURITY-TOKEN"),
			CST:              response.Header.Get("CST"),
		},
	}

	return newSessionResponse, headerTokens, nil
}

func (capClient *CapitalClientAPI) GetEncryptedPassword(encriptionResponse EncriptionResponse) (string, error) {
	input := []byte(capClient.CAPITAL_API_KEY_PASSWORD + "|" + strconv.FormatInt(int64(encriptionResponse.TimeStamp), 10))
	input = []byte(base64.StdEncoding.EncodeToString(input))
	publicKey, err := base64.StdEncoding.DecodeString(encriptionResponse.EncryptionKey)
	if err != nil {
		return "", err
	}
	pubKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("Not a valid RSA public key")
	}
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, input)
	if err != nil {
		return "", err
	}
	output := base64.StdEncoding.EncodeToString(cipher)

	return output, nil
}

func (capClient *CapitalClientAPI) GetAllAccounts() (AccountsResponse, error) {
	if capClient.HttpClient.Transport == nil {
		return AccountsResponse{}, errors.New("A Session is needed; Run `capClient.CreateNewSession()` to authenticate")
	}
	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/accounts", nil)
	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return AccountsResponse{}, err
	}
	defer response.Body.Close()

	var accountsResponse AccountsResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&accountsResponse)

	return accountsResponse, nil
}

func (capClient *CapitalClientAPI) GetMarketsDetails(epics []string) (MarketsDetailsResponse, error) {
	if capClient.HttpClient.Transport == nil {
		return MarketsDetailsResponse{}, errors.New("A session is need; Run `capClient.CreateNewSession()` to authenticate first")
	}

	values := url.Values{}
	//if searchTerm != "" {
	//	values.Set("searchTerm", searchTerm)
	//}
	for _, epic := range epics {
		values.Add("epics", epic)
	}

	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/markets", nil)
	request.URL.RawQuery = values.Encode()
	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return MarketsDetailsResponse{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return MarketsDetailsResponse{}, errors.New(fmt.Sprintf("Unexpected [%d] Status Code Response - %s", response.StatusCode, string(body)))
	}

	var marketsResponse MarketsDetailsResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&marketsResponse)

	return marketsResponse, nil
}

func (capClient *CapitalClientAPI) GetPrices(epic string, resolution string) (pricesResponse PricesResponse, err error) {
	if capClient.HttpClient.Transport == nil {
		return pricesResponse, errors.New("A session is need; Run `capClient.CreateNewSession()` to authenticate first")
	}

	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/prices/"+epic, nil)
	values := request.URL.Query()
	values.Set("max", "100")
	values.Set("resolution", resolution)
	request.URL.RawQuery = values.Encode()

	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return pricesResponse, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return pricesResponse, errors.New(fmt.Sprintf("Unexpected [%d] Status Code Response - %s", response.StatusCode, string(body)))
	}
	pricesResponse = PricesResponse{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&pricesResponse)

	return pricesResponse, nil
}

func (capClient *CapitalClientAPI) GetPositions() (positionsResponse PositionsResponse, err error) {
	if capClient.HttpClient.Transport == nil {
		return positionsResponse, errors.New("A session is need; Run `capClient.CreateNewSession()` to authenticate first")
	}

	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/positions", nil)
	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return positionsResponse, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return positionsResponse, errors.New(fmt.Sprintf("Unexpected [%d] Status Code Response - %s", response.StatusCode, string(body)))
	}
	positionsResponse = PositionsResponse{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&positionsResponse)

	return positionsResponse, nil
}

type Direction string

const (
	SELL Direction = "SELL"
	BUY  Direction = "BUY"
)

func (capClient *CapitalClientAPI) CreatePosition(epic string, direction Direction, size float64) (createPositionResponse CreatePositionResponse, err error) {
	if capClient.HttpClient.Transport == nil {
		return createPositionResponse, errors.New("A session is need; Run `capClient.CreateNewSession()` to authenticate first")
	}

	body, err := json.Marshal(CreatePositionBody{
		Epic:      epic,
		Direction: direction,
		Size:      fmt.Sprintf("%v", size),
	})
	if err != nil {
		return createPositionResponse, err
	}

	request, _ := http.NewRequest("POST", "https://api-capital.backend-capital.com/api/v1/positions", bytes.NewBuffer(body))
	response, err := capClient.HttpClient.Do(request)
	if err != nil {
		return createPositionResponse, err
	}
	defer response.Body.Close()

	createPositionResponse = CreatePositionResponse{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&createPositionResponse)

	return createPositionResponse, nil
}
