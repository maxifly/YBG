package main

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
	"os"
)

func GetCheckCodeUrl(clientId string) string {
	return "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + clientId
}

func CreateToken(clientId string, clientSecret string, code string) (TokenInfo, error) {

	config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     yandex.Endpoint,
	}

	tokenValue, err := config.Exchange(context.Background(), code)
	if err != nil {
		return *new(TokenInfo), nil
	}
	tokenInfo := TokenInfo{AccessToken: tokenValue.AccessToken,
		RefreshToken: tokenValue.RefreshToken,
		Expiry:       tokenValue.Expiry}
	return tokenInfo, err
}

func WriteToken(tokenInfo TokenInfo) error {
	jsonData, err := json.Marshal(tokenInfo)
	if err != nil {
		return err
	}

	err = os.WriteFile(FILE_PATH_TOKEN, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
