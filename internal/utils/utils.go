package utils

import (
	"os"
	"path/filepath"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckForConfig() (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	filename := filepath.Join(cwd, "skelly.yaml")

	_, err = os.Stat(filename)
	if err == nil {
		// file exists
		return true, nil
	}
	if os.IsNotExist(err) {
		// file does _not_ exist
		return false, nil
	}
	return false, err
}
