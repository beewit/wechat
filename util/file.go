package util

import (
	"strings"
	"os"
)

func MkDirAll(path string) error {
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = strings.Replace(path, "/", "\\", -1)
	}
	flog, err := PathExists(path)
	if err != nil {
		return err
	}
	if flog {
		return nil
	}
	err2 := os.MkdirAll(path, os.ModePerm)

	if err2 != nil {
		return err2
	}
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
