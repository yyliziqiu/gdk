package xfile

import (
	"fmt"
	"os"
)

// Exist 判断文件是否存在
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MakeDir 创建目录，如果目录已存在则不做任何操作
func MakeDir(path string) error {
	stat, err := os.Stat(path)
	if err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("%s is not a directory", path)
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, 0755)
}
