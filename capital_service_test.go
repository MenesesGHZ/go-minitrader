package forexbot

import (
	"fmt"
	"testing"
)

func TestGetEncriptionKey(t *testing.T) {
	encription, err := GetEncriptionKey()
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
	encription, _ := GetEncriptionKey()
	encryptedPassword, err := GetEncryptedPassword(encription)
	if err != nil {
		fmt.Println(err)
		t.Errorf("Got Error")
	}
	t.Logf("Gathered Key: %s", encryptedPassword)
}

func TestCreateNewSessionAccount(t *testing.T) {
	newSessionResponse, headerTokens, err := CreateNewSession()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Got Error")
	}
	t.Logf("SessionResponse: %+v", newSessionResponse)
	t.Logf("Token Headers: %+v", headerTokens)
}
