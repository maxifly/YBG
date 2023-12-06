package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type BackupFileInfo struct {
	Name       string
	CreateDate string
	Size       string
	IsLocal    bool
	IsRemote   bool
}

type LocalBackupFile struct {
	Slug       string
	Name       string
	CreateDate string
	Path       string
	Size       string
}

func GetFilesInfo(application *Application) ([]BackupFileInfo, error) {
	application.debugLog.Println("Start get files")

	getRemoteFiles(application)

	//TODO Create
	bFiles := make([]BackupFileInfo, 0, 0)

	bFiles = append(bFiles, BackupFileInfo{Name: "test1", CreateDate: "01.01", Size: "125", IsLocal: true, IsRemote: true})
	bFiles = append(bFiles, BackupFileInfo{Name: "test2", CreateDate: "02.02", Size: "125", IsLocal: true, IsRemote: true})

	return bFiles, nil
}

func getLocalBacupFiles(app *Application) (error, map[string] LocalBackupFile) {

	entries, err := os.ReadDir(BACKUP_PATH)
	if err != nil {
		app.errorLog.Printf("Unable to read backup %s. %v", BACKUP_PATH, err)
		return fmt.Errorf("error when read local backups"), nil
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil { ... }
		infos = append(infos, info)
	}

	//files, err := ioutil.ReadDir(".")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, file := range files {
	//	fmt.Println(file.Name(), file.IsDir())
	//}
}

func extract(tarfile string) {
	reader, err := os.Open(tarfile)
	if err != nil {
		fmt.Printf("ERROR: cannot read tar file, error=[%v]\n", err)
		return
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("ERROR: cannot read tar file, error=[%v]\n", err)
			return
		}

		j, err := json.Marshal(header)
		if err != nil {
			fmt.Printf("ERROR: cannot parse header, error=[%v]\n", err)
			return
		}
		fmt.Printf("header=%s\n", string(j))

		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(header.Name, 0755); err != nil {
				fmt.Printf("ERROR: cannot mkdir file, error=[%v]\n", err)
				return
			}
		} else {
			if err = os.MkdirAll(path.Dir(header.Name), 0755); err != nil {
				fmt.Printf("ERROR: cannot file mkdir file, error=[%v]\n", err)
				return
			}

			file, err := os.OpenFile(header.Name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				fmt.Printf("ERROR: cannot open file, error=[%v]\n", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, tarReader)
			if err != nil {
				fmt.Printf("ERROR: cannot write file, error=[%v]\n", err)
				return
			}
		}
	}
}
