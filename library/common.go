package library

import (
	"os"
	"path/filepath"
)

func Log() {

}

/**
文件(夹)路径是否存在
*/
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

/**
获取当前路径
*/
func GetCurrentDirectory() (path string, err error) {
	path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}
