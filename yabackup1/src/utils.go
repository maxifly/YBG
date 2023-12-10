package main

import (
	"fmt"
	"time"
)

const MiB = float64(1 << 20)

func (fs fileSize) Convert2MbString() string {
	return fmt.Sprintf("%.2f", float64(fs)/MiB)
}

func (fm fileModified) Convert2String() string {
	return time.Time(fm).Format("02.01.2006 15:04:05 MST")
}
