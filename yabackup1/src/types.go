package main

import "time"

type fileSize int64
type fileModified time.Time

type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

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
