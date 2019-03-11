package util

import "time"

func GetCurrTimeStamp() int64 {
	return time.Now().Unix()
}
