package emojifilter

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// isExist check path is exists, exist return true, not exist return false
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// Check if error is "no such file or directory"
	if _, ok := err.(*os.PathError); ok {
		return false
	}
	return false
}

// readFile read file content
func readFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// writeFile read content to the file
func writeFile(file string, body []byte) (int, error) {
	err := mkdirAll(filepath.Dir(file))
	if err != nil {
		return 0, err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	n, err := f.Write(body)
	f.Close()
	return n, err
}

// mkdirAll check the path isexist or mkdir, and check writable
func mkdirAll(path string) error {
	var err error
	// check path exist or create
	if isExist(path) == false {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	// check path writable
	if isWritable(path) == false {
		return errors.New("path [" + path + "] is not writable!")
	}
	return nil
}

// isWritable check path is writeable, can return true, can not return false
func isWritable(path string) bool {
	err := unix.Access(path, unix.O_RDWR)
	if err == nil {
		return true
	}
	// Check if error is "no such file or directory"
	if _, ok := err.(*os.PathError); ok {
		return false
	}
	return false
}

// httpDownload HTTP下载文件
func httpDownload(uri, path string) (int, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return writeFile(path, body)
}
