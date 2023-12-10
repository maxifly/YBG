package main

import "time"

type fileSize int64
type fileModified time.Time

type GeneralFileInfo struct {
	Name     string
	Size     fileSize
	Modified fileModified
}

type RemoteFileInfo GeneralFileInfo

type BackupFileInfo struct {
	GeneralInfo GeneralFileInfo
	IsLocal     bool
	IsRemote    bool
}

type LocalBackupFileInfo struct {
	GeneralInfo GeneralFileInfo
	Slug        string
	Path        string
}

type BackupArchInfo struct {
	Slug string
	Name string
}
