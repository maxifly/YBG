package main

import (
	"context"
	yadisk "github.com/maxifly/yandex-disk-sdk-go"
	"net/http"
)

const itemTypeFile string = "file"

type RemoteFileInfo struct {
	Name     string
	Size     int
	Modified string
}

func NewYandexDisk(accessToken string) (yadisk.YaDisk, error) {
	return yadisk.NewYaDisk(context.Background(), http.DefaultClient, &yadisk.Token{AccessToken: accessToken})

}

func getRemoteFiles(app *Application) []RemoteFileInfo {
	disk, err := (*app.yaDisk).GetDisk([]string{})
	app.infoLog.Printf("%v", err)
	app.infoLog.Printf("%v", disk)
	app.infoLog.Printf("%v", app.options.RemotePath)

	result := make([]RemoteFileInfo, 0)
	resource, err := (*app.yaDisk).GetResource(app.options.RemotePath, make([]string, 0), 0, 10000, false, "0", "name")
	if err != nil {
		app.errorLog.Printf("Disk %v", (app.yaDisk))
		app.errorLog.Printf("Error when get remote files %v", err)
		return result
	}

	app.infoLog.Printf("%v", resource)

	for _, item := range resource.Embedded.Items {
		if item.Type != itemTypeFile {
			continue
		}
		result = append(result, RemoteFileInfo{Name: item.Name,
			Size:     item.Size,
			Modified: item.Modified})

	}
	return result

}
