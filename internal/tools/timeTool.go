package tools

import (
	"time"
)

const (
	timeLoc = "Asia/Shanghai"
)

var loc, _ = time.LoadLocation(timeLoc)

func GetNowTime() time.Time {
	shanghaiTime := time.Now().In(loc)
	return shanghaiTime
}
