package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func GetFilesInfo(application *Application) ([]BackupFileInfo, error) {
	application.debugLog.Println("Start get files")

	remoteFiles := getRemoteFiles(application)
	err, localFiles := getLocalBackupFiles(application)
	if err != nil {
		localFiles = make(map[string]LocalBackupFileInfo)
	}

	return intersectFiles(application, localFiles, remoteFiles)
}

type stringSet map[string]bool

func intersectFiles(app *Application,
	localFiles map[string]LocalBackupFileInfo,
	remoteFiles []RemoteFileInfo) ([]BackupFileInfo, error) {

	remoteFileNames := make(stringSet)
	processedRemoteFile := make(stringSet)

	for _, remoteFile := range remoteFiles {
		remoteFileNames[remoteFile.Name] = true
	}

	result := make([]BackupFileInfo, 0, len(localFiles))

	// Обработаем локальные файлы
	for _, localFile := range localFiles {
		remoteFileName := generateRemoteFileName(localFile)
		_, isRemote := remoteFileNames[remoteFileName]

		result = append(result,
			BackupFileInfo{
				GeneralInfo: localFile.GeneralInfo,
				IsLocal:     true,
				IsRemote:    isRemote,
			})

		processedRemoteFile[remoteFileName] = true
	}

	for _, remoteFile := range remoteFiles {
		if _, isProcessing := processedRemoteFile[remoteFile.Name]; !isProcessing {
			result = append(result,
				BackupFileInfo{
					GeneralInfo: GeneralFileInfo{
						Name:     remoteFile.Name,
						Size:     remoteFile.Size,
						Modified: remoteFile.Modified,
					},
					IsLocal:  false,
					IsRemote: true,
				})
			processedRemoteFile[remoteFile.Name] = true
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return time.Time(result[i].GeneralInfo.Modified).After(time.Time(result[j].GeneralInfo.Modified))
	})
	return result, nil
}

// localFile.CreateDate.Format("02.01.2006 15:04:05 MST"),
func generateRemoteFileName(localFile LocalBackupFileInfo) string {
	return strings.ReplaceAll(strings.ReplaceAll((localFile.GeneralInfo.Name+"_"+localFile.Slug), " ", "-"), ":", "_")
}

func getLocalBackupFiles(app *Application) (error, map[string]LocalBackupFileInfo) {

	entries, err := os.ReadDir(BACKUP_PATH)
	if err != nil {
		app.errorLog.Printf("Unable to read backup %s. %v", BACKUP_PATH, err)
		return fmt.Errorf("error when read local backups"), nil
	}
	result := make(map[string]LocalBackupFileInfo)
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

		result[archInfo.Slug] = LocalBackupFileInfo{
			GeneralInfo: convertYdInfoToGeneral(&info),
			Slug:        archInfo.Slug,
			Path:        filePath,
		}
	}
	return nil, result

}

func convertYdInfoToGeneral(ydFileInfo *fs.FileInfo) GeneralFileInfo {
	return GeneralFileInfo{Name: (*ydFileInfo).Name(),
		Size:     fileSize((*ydFileInfo).Size()),
		Modified: fileModified((*ydFileInfo).ModTime()),
	}
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

			app.infoLog.Printf("plan=%v\n", plan)
			err = json.Unmarshal(plan, &data)
			app.infoLog.Printf("data=%v\n", data)
			if err != nil {
				return nil, fmt.Errorf("cannot parse backup info, error=[%v]", err)
			}
			if data.Slug == "" || data.Name == "" {
				return nil, fmt.Errorf("cannot parse backup info. Necessary field not found")
			}
			return &data, nil
		}
	}
	return nil, fmt.Errorf("backup info not found")
}
