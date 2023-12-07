package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	Size       int64
}

func GetFilesInfo(application *Application) ([]BackupFileInfo, error) {
	application.debugLog.Println("Start get files")

	getRemoteFiles(application)
	err, localFiles := getLocalBackupFiles(application)
	if err != nil {
		return nil, err
	}

	application.debugLog.Printf("Local files: %+v", localFiles)
	//TODO Create
	bFiles := make([]BackupFileInfo, 0, 0)

	bFiles = append(bFiles, BackupFileInfo{Name: "test1", CreateDate: "01.01", Size: "125", IsLocal: true, IsRemote: true})
	bFiles = append(bFiles, BackupFileInfo{Name: "test2", CreateDate: "02.02", Size: "125", IsLocal: true, IsRemote: true})

	return bFiles, nil
}

func getLocalBackupFiles(app *Application) (error, map[string]LocalBackupFile) {

	entries, err := os.ReadDir(BACKUP_PATH)
	if err != nil {
		app.errorLog.Printf("Unable to read backup %s. %v", BACKUP_PATH, err)
		return fmt.Errorf("error when read local backups"), nil
	}
	result := make(map[string]LocalBackupFile)
	for _, entry := range entries {
		app.debugLog.Printf("entry %+v", entry)
		info, err := entry.Info()
		if err != nil {
			app.errorLog.Printf("Error read file info %v", err)
			continue
		}
		app.debugLog.Printf("info: %+v", info)

		if info.IsDir() {
			continue
		}

		filePath := filepath.Join(BACKUP_PATH, info.Name())
		app.debugLog.Printf("Read %s", filePath)
		archInfo, err := extractArchInfo(app, filePath)
		if err != nil {
			app.errorLog.Printf("Error extract slug from %s %v", info.Name(), err)
			continue
		}

		result[archInfo.Slug] = LocalBackupFile{Slug: archInfo.Slug,
			Name:       info.Name(),
			Path:       filePath,
			Size:       info.Size(),
			CreateDate: info.ModTime().String(),
		}
	}
	return nil, result

	//files, err := ioutil.ReadDir(".")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, file := range files {
	//	fmt.Println(file.Name(), file.IsDir())
	//}
}

type BackupArchInfo struct {
	Slug string
	Name string
}

func extractArchInfo(app *Application, tarfile string) (*BackupArchInfo, error) {
	reader, err := os.Open(tarfile)
	if err != nil {
		return nil, fmt.Errorf("ERROR: cannot read tar file, error=[%v]\n", err)
	}

	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			app.errorLog.Printf("Can not close reader, error=[%v]", err)
		}
	}(reader)

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("cannot read tar file, error=[%v]", err)
		}

		j, err := json.Marshal(header)
		if err != nil {
			return nil, fmt.Errorf("cannot parse header, error=[%v]", err)
		}
		app.infoLog.Printf("header=%s\n", string(j))

		info := header.FileInfo()
		if info.IsDir() || info.Name() != "backup.json" {
			continue
		} else {
			var data BackupArchInfo
			plan, err := io.ReadAll(tarReader)

			if err != nil {
				return nil, fmt.Errorf("cannot read backup info, error=[%v]", err)

			}

			err = json.Unmarshal(plan, &data)
			if err != nil {
				return nil, fmt.Errorf("cannot parse backup info, error=[%v]", err)

			}
			return &data, nil
		}
	}
	return nil, fmt.Errorf("backup info not found")
}
