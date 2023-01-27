package forexbot

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const CAPITAL_EMAIL = ""
const CAPITAL_API_KEY = ""
const CAPITAL_API_KEY_PASSWORD = ""

func GetEncriptionKey() (EncriptionResponse, error) {
	client := &http.Client{}

	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/session/encryptionKey", nil)
	request.Header.Add("X-CAP-API-KEY", CAPITAL_API_KEY)

	response, err := client.Do(request)
	if err != nil {
		return EncriptionResponse{}, err
	}
	defer response.Body.Close()

	var encription EncriptionResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&encription)

	return encription, nil
}

func GetWatchLists() (WatchListsResponse, error) {
	client := &http.Client{}

	request, _ := http.NewRequest("GET", "https://api-capital.backend-capital.com/api/v1/watchlists", nil)
	request.Header.Add("X-CAP-API-KEY", CAPITAL_API_KEY)

	response, err := client.Do(request)
	if err != nil {
		return WatchListsResponse{}, err
	}
	defer response.Body.Close()

	var watchlists WatchListsResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&watchlists)

	return watchlists, nil
}

func CreateNewSession() (newSessionResponse NewSessionResponse, headerTokens http.Header, err error) {
	encriptionResponse, err := GetEncriptionKey()
	if err != nil {
		return newSessionResponse, headerTokens, err
	}
	encriptedPassword, err := GetEncryptedPassword(encriptionResponse)
	if err != nil {
		return newSessionResponse, headerTokens, err
	}

	body := NewSessionBody{
		Identifier:        CAPITAL_EMAIL,
		Password:          encriptedPassword,
		EncryptedPassword: true,
	}

	jsonData, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "https://api-capital.backend-capital.com/api/v1/watchlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-CAP-API-KEY", CAPITAL_API_KEY)

	client := &http.Client{}
	response, err := client.Do(request)
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
	headerTokens.Add("CST", response.Header.Get("CST"))
	headerTokens.Add("X-SECURITY-TOKEN", response.Header.Get("X-SECURITY-TOKEN"))

	return newSessionResponse, headerTokens, nil
}

func GetEncryptedPassword(encription EncriptionResponse) (string, error) {
	input := []byte(CAPITAL_API_KEY_PASSWORD + "|" + strconv.FormatInt(int64(encription.TimeStamp), 10))
	publicKey, err := base64.StdEncoding.DecodeString(encription.EncryptionKey)
	if err != nil {
		return "", err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	rsaKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("key is not a RSA key")
	}
	encodedInput := base64.StdEncoding.EncodeToString(input)
	hashed := sha256.Sum256([]byte(encodedInput))
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, hashed[:], nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
