package main

import (
	"context"
	yadisk "github.com/nikitaksv/yandex-disk-sdk-go"
	"net/http"
)

func NewYandexDisk(accessToken string) (yadisk.YaDisk, error) {
	return yadisk.NewYaDisk(context.Background(), http.DefaultClient, &yadisk.Token{AccessToken: accessToken})
}

//func get_files(app *Application) error {
//
//}
