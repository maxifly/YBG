package main

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

func GetCheckCodeUrl(clientId string) string {
	return "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + clientId
}

func CreateToken(clientId string, clientSecret string, code string) (oauth2.Token, error) {

	config := oauth2.Config{
		ClientID:     "859eee7fe42742f485c982b813646431",
		ClientSecret: "4434d58aa29141529ddd9fd6193caf4e",
		Endpoint:     yandex.Endpoint,
	}

	tokenValue, err := config.Exchange(context.Background(), code)
	if err == nil {
		return *tokenValue, nil
	}
	return oauth2.Token{}, err

}
