package util

import (
	"log"
)

var LogInfoln = log.Println

// 沒有加锁, 请确保在初始化阶段调用!
func SetLogInfoln(fn func(v ...interface{})) {
	if fn == nil {
		return
	}
	LogInfoln = fn
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
