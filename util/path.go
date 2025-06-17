package util

import "os"

func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FolderIsNotEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Readdir(1)
	return err == nil
}
