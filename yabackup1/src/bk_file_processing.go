package main

type BackupFileInfo struct {
	Name       string
	CreateDate string
	Size       string
	IsLocal    bool
	IsRemote   bool
}

func GetFilesInfo(application *Application) ([]BackupFileInfo, error) {
	application.debugLog.Println("Start get files")
	//TODO Create
	bFiles := make([]BackupFileInfo, 0, 0)

	bFiles = append(bFiles, BackupFileInfo{Name: "test1", CreateDate: "01.01", Size: "125", IsLocal: true, IsRemote: true})
	bFiles = append(bFiles, BackupFileInfo{Name: "test2", CreateDate: "02.02", Size: "125", IsLocal: true, IsRemote: true})

	return bFiles, nil
}
