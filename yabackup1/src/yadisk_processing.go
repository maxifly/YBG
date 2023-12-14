package main

import (
	"context"
	"fmt"
	yadisk "github.com/nikitaksv/yandex-disk-sdk-go"
	"net/http"
	"time"
)

const itemTypeFile string = "file"

func NewYandexDisk(accessToken string) (yadisk.YaDisk, error) {
	return yadisk.NewYaDisk(context.Background(), http.DefaultClient, &yadisk.Token{AccessToken: accessToken})

}

func getRemoteFiles(app *Application) ([]RemoteFileInfo, error) {
	app.infoLog.Printf("%v", app.options.RemotePath)

	result := make([]RemoteFileInfo, 0)
	resource, err := (*app.yaDisk).GetResource(app.options.RemotePath, make([]string, 0), 10000, 0, false, "0", "name")
	if err != nil {
		app.errorLog.Printf("Error when get remote files %v", err)
		return result, err
	}

	app.debugLog.Printf("Found %d items", len(resource.Embedded.Items))

	for _, item := range resource.Embedded.Items {
		if item.Type != itemTypeFile {
			continue
		}

		modifyedTime, err := convertDateString(item.Modified)
		if err != nil {
			app.errorLog.Printf("Can not parse data %s %v", item.Modified, err)
			modifyedTime = time.Now() //TODO Сделат какую-то минимальную дату по умолчанию
		}
		result = append(result, RemoteFileInfo{Name: item.Name,
			Size:     fileSize(item.Size),
			Modified: fileModified(modifyedTime)})

	}

	app.debugLog.Printf("Processing %d files", len(result))
	return result, nil

}

func uploadFile(source string, destination string) error {
	//TODO create
	return fmt.Errorf("error upload %s %s", source, destination)
}
func convertDateString(modified string) (time.Time, error) {
	return time.Parse(time.RFC3339, modified)
	//"modified": "2023-10-31T03:32:52+00:00",
}
